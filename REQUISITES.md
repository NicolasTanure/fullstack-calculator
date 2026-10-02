# Calculator Requirements

## Objective

Build a full-stack calculator with a React frontend and a Go REST API.
Focus on correct behavior, clear design, maintainable code, and useful tests.

## Functional requirements

| Area | Requirements |
| --- | --- |
| Operations | Addition, subtraction, multiplication, division |
| Frontend | Two operands, operation selection, clear results, input validation, error handling, basic mobile support |
| Backend | Operation endpoints, input validation, relevant edge cases, division-by-zero handling, JSON results and errors |

Exponentiation, square root, and percentage are optional and come after the
required scope is complete. Their behavior must be agreed before implementation.

## Technologies and delivery

- React and TypeScript frontend; Go backend.
- Clean, readable, idiomatic code with separate responsibilities.
- Unit tests covering important behavior in both applications, plus coverage reports.
- Dockerfiles and Docker Compose to run both applications together.
- README with setup instructions, how to run both applications, API examples,
  design decisions, and assumptions.
- A Git repository containing the source and documentation. Publication requires
  the user's explicit authorization.
- Existing AI development notes in `PROMPTS.md`; add new prompts only when the
  user explicitly asks.

Prioritize correctness, clarity, and maintainability over extra features.
`TASK.md` defines the workflow; `AGENTS.md` records the approved implementation details.
