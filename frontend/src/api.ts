const invalidResponse = 'The server returned an invalid response. Please try again.'
const unavailable = 'The calculator service is unavailable. Please try again.'
const timeoutMessage = 'The request timed out. Please try again.'

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export type Operation = 'add' | 'subtract' | 'multiply' | 'divide'

export async function calculate(operation: Operation, left: number, right: number): Promise<number> {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 10_000)

  try {
    let response: Response
    try {
      response = await fetch(`/api/${operation}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ left, right }),
        signal: controller.signal,
      })
    } catch {
      throw new Error(controller.signal.aborted ? timeoutMessage : unavailable)
    }

    let body: unknown
    try {
      const contentType = response.headers.get('Content-Type')?.split(';')[0].trim().toLowerCase()
      if (contentType !== 'application/json') throw new Error(invalidResponse)
      body = await response.json()
    } catch {
      if (controller.signal.aborted) throw new Error(timeoutMessage)
      throw new Error(response.status >= 500 ? unavailable : invalidResponse)
    }

    if (response.ok) {
      if (!isObject(body) || typeof body.result !== 'number' || !Number.isFinite(body.result)) {
        throw new Error(invalidResponse)
      }
      return body.result
    }

    if (
      isObject(body) && isObject(body.error) &&
      typeof body.error.code === 'string' && body.error.code.trim() !== '' &&
      typeof body.error.message === 'string' && body.error.message.trim() !== ''
    ) {
      throw new Error(body.error.message)
    }
    throw new Error(response.status >= 500 ? unavailable : invalidResponse)
  } finally {
    clearTimeout(timeout)
  }
}
