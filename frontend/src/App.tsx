import { useEffect, useRef, useState, type FormEvent } from 'react'
import { calculate as calculateRequest, type Operation } from './api'
import { formatResult, needsApproximationNotice, parseOperand, type OperandError } from './numbers'

function operandErrorMessage(field: string, error: OperandError): string {
  if (error === 'out_of_range') {
    return `${field} is outside the supported numeric range. Enter a smaller absolute value.`
  }
  return `${field} has an invalid format. Use a number such as -2.5 or 1e3.`
}

type ErrorFeedback = {
  message: string
  field?: 'left' | 'right'
}

type CalculationResult = {
  value: number
  approximate: boolean
}

export default function App() {
  const leftInput = useRef<HTMLInputElement>(null)
  const rightInput = useRef<HTMLInputElement>(null)
  const calculateButton = useRef<HTMLButtonElement>(null)
  const restoreRequestFocus = useRef(false)
  const [left, setLeft] = useState('')
  const [right, setRight] = useState('')
  const [operation, setOperation] = useState<Operation>('add')
  const [result, setResult] = useState<CalculationResult | null>(null)
  const [error, setError] = useState<ErrorFeedback | null>(null)
  const [pending, setPending] = useState(false)

  useEffect(() => {
    if (!pending && restoreRequestFocus.current) {
      restoreRequestFocus.current = false
      calculateButton.current?.focus()
    }
  }, [pending])

  function clearOutput() {
    setResult(null)
    setError(null)
  }

  async function calculate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (pending) return
    clearOutput()

    const leftNumber = parseOperand(left)
    const rightNumber = parseOperand(right)
    if (!leftNumber.valid) {
      setError({ message: operandErrorMessage('First number', leftNumber.error), field: 'left' })
      leftInput.current?.focus()
      return
    }
    if (!rightNumber.valid) {
      setError({ message: operandErrorMessage('Second number', rightNumber.error), field: 'right' })
      rightInput.current?.focus()
      return
    }

    restoreRequestFocus.current = true
    setPending(true)
    try {
      const value = await calculateRequest(operation, leftNumber.value, rightNumber.value)
      setResult({
        value,
        approximate: needsApproximationNotice(value, [leftNumber.value, rightNumber.value]),
      })
    } catch (failure) {
      setError({
        message: failure instanceof Error ? failure.message : 'Unable to calculate. Please try again.',
      })
    } finally {
      setPending(false)
    }
  }

  return (
    <main className="calculator-shell">
      <section className="calculator-card" aria-labelledby="calculator-title">
        <h1 id="calculator-title">Calculator</h1>
        <p className="introduction">Simple arithmetic, clear results.</p>
        <form onSubmit={calculate} aria-busy={pending} noValidate>
          <div className="operands">
            <div className="field">
              <label htmlFor="left">First number</label>
              <input
                id="left"
                ref={leftInput}
                type="text"
                inputMode="decimal"
                aria-invalid={error?.field === 'left'}
                aria-describedby={error?.field === 'left' ? 'number-format calculation-error' : 'number-format'}
                value={left}
                disabled={pending}
                onChange={(event) => {
                  setLeft(event.target.value)
                  clearOutput()
                }}
              />
            </div>
            <div className="field">
              <label htmlFor="right">Second number</label>
              <input
                id="right"
                ref={rightInput}
                type="text"
                inputMode="decimal"
                aria-invalid={error?.field === 'right'}
                aria-describedby={error?.field === 'right' ? 'number-format calculation-error' : 'number-format'}
                value={right}
                disabled={pending}
                onChange={(event) => {
                  setRight(event.target.value)
                  clearOutput()
                }}
              />
            </div>
          </div>
          <p id="number-format" className="input-hint">
            Use a period for decimals. Scientific notation such as 1e3 is supported.
          </p>
          <div className="field">
            <label htmlFor="operation">Operation</label>
            <select
              id="operation"
              disabled={pending}
              value={operation}
              onChange={(event) => {
                setOperation(event.target.value as Operation)
                clearOutput()
              }}
            >
              <option value="add">Addition (+)</option>
              <option value="subtract">Subtraction (−)</option>
              <option value="multiply">Multiplication (×)</option>
              <option value="divide">Division (÷)</option>
            </select>
          </div>
          <button ref={calculateButton} type="submit" disabled={pending}>
            {pending ? 'Calculating…' : 'Calculate'}
          </button>
        </form>
        {error !== null && <p id="calculation-error" className="error" role="alert">{error.message}</p>}
        {result !== null && (
          <div className="result">
            <span className="result-label">Result</span>
            <output
              aria-label="Result"
              aria-live="polite"
              aria-describedby={result.approximate ? 'precision-note' : undefined}
            >
              {result.approximate ? '≈ ' : ''}{formatResult(result.value)}
            </output>
            {result.approximate && (
              <p id="precision-note" className="result-hint">
                Approximate result. Large numbers can lose precision; results are shown with up to 12 significant digits.
              </p>
            )}
          </div>
        )}
      </section>
    </main>
  )
}
