import { afterEach, beforeEach, expect, test, vi } from 'vitest'
import { calculate } from './api'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.useRealTimers()
})

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

test('sends numeric operands to the addition endpoint and preserves the API result', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 0.30000000000000004 }))

  await expect(calculate('add', 0.1, 0.2)).resolves.toBe(0.30000000000000004)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: '{"left":0.1,"right":0.2}',
    signal: expect.any(AbortSignal),
  })
})

test('propagates a valid API error message', async () => {
  fetchMock.mockResolvedValue(jsonResponse({
    error: { code: 'non_finite_result', message: 'The calculation result must be finite.' },
  }, 400))
  await expect(calculate('add', 1e308, 1e308)).rejects.toThrow('The calculation result must be finite.')
})

test('sends subtraction operands in their original order', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: -6 }))
  await expect(calculate('subtract', 2, 8)).resolves.toBe(-6)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/subtract', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: '{"left":2,"right":8}',
    signal: expect.any(AbortSignal),
  })
})

test('sends numeric operands to the multiplication endpoint', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 0.020000000000000004 }))
  await expect(calculate('multiply', 0.1, 0.2)).resolves.toBe(0.020000000000000004)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/multiply', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: '{"left":0.1,"right":0.2}',
    signal: expect.any(AbortSignal),
  })
})

test('sends division operands in their original order', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 0.25 }))
  await expect(calculate('divide', 2, 8)).resolves.toBe(0.25)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/divide', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: '{"left":2,"right":8}',
    signal: expect.any(AbortSignal),
  })
})

test.each([null, [], {}, { result: null }, { result: '10' }])(
  'rejects an invalid success response %j', async (body) => {
    fetchMock.mockResolvedValue(jsonResponse(body))
    await expect(calculate('add', 8, 2)).rejects.toThrow('The server returned an invalid response.')
  },
)

test('rejects a JSON result outside the finite range', async () => {
  fetchMock.mockResolvedValue(new Response('{"result":1e309}', {
    headers: { 'Content-Type': 'application/json' },
  }))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The server returned an invalid response.')
})

test('rejects malformed JSON', async () => {
  fetchMock.mockResolvedValue(new Response('{', {
    headers: { 'Content-Type': 'application/json' },
  }))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The server returned an invalid response.')
})

test('rejects an incompatible response content type', async () => {
  fetchMock.mockResolvedValue(new Response('{"result":10}', {
    headers: { 'Content-Type': 'text/plain' },
  }))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The server returned an invalid response.')
})

test.each([
  { error: { code: 400, message: 'Rejected.' } },
  { error: { code: 'invalid_input', message: 123 } },
  { error: { code: 'invalid_input', message: ' ' } },
  { error: { code: ' ', message: 'Rejected.' } },
])('rejects an invalid error envelope %j', async (body) => {
  fetchMock.mockResolvedValue(jsonResponse(body, 400))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The server returned an invalid response.')
})

test('reports unavailability for an invalid JSON error envelope from a failing service', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ detail: 'Upstream failed.' }, 503))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The calculator service is unavailable.')
  expect(fetchMock).toHaveBeenCalledTimes(1)
})

test('reports service unavailability for a proxy failure', async () => {
  fetchMock.mockResolvedValue(new Response('', { status: 502 }))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The calculator service is unavailable.')
})

test('reports a network failure without automatic retries', async () => {
  fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
  await expect(calculate('add', 8, 2)).rejects.toThrow('The calculator service is unavailable.')
  expect(fetchMock).toHaveBeenCalledTimes(1)
})

test('aborts the request at 10 seconds and clears its timer', async () => {
  vi.useFakeTimers()
  fetchMock.mockImplementation((_, options) => new Promise((_, reject) => {
    options?.signal?.addEventListener('abort', () => {
      reject(new DOMException('Aborted', 'AbortError'))
    }, { once: true })
  }))

  const failure = expect(calculate('add', 8, 2)).rejects.toThrow('The request timed out.')
  await vi.advanceTimersByTimeAsync(9_999)
  expect(fetchMock.mock.calls[0][1]?.signal?.aborted).toBe(false)
  await vi.advanceTimersByTimeAsync(1)
  await failure
  expect(fetchMock.mock.calls[0][1]?.signal?.aborted).toBe(true)
  expect(vi.getTimerCount()).toBe(0)
  expect(fetchMock).toHaveBeenCalledTimes(1)
})

test('clears the timeout after a successful request', async () => {
  vi.useFakeTimers()
  fetchMock.mockResolvedValue(jsonResponse({ result: 10 }))
  await expect(calculate('add', 8, 2)).resolves.toBe(10)
  expect(vi.getTimerCount()).toBe(0)
})

test('keeps the deadline active while reading the response body', async () => {
  vi.useFakeTimers()
  fetchMock.mockImplementation(async (_, options) => {
    const response = jsonResponse({ result: 10 })
    vi.spyOn(response, 'json').mockImplementation(() => new Promise((_, reject) => {
      options?.signal?.addEventListener('abort', () => {
        reject(new DOMException('Aborted', 'AbortError'))
      }, { once: true })
    }))
    return response
  })
  const failure = expect(calculate('add', 8, 2)).rejects.toThrow('The request timed out.')
  await vi.advanceTimersByTimeAsync(10_000)
  await failure
  expect(vi.getTimerCount()).toBe(0)
})
