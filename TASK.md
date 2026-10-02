# Calculator Project

Build a small full-stack calculator incrementally, keeping the code clear,
well structured, and easy to explain.

## Objective and priorities

Use React and TypeScript for the frontend and Go for the backend. The frontend
must call the REST API for arithmetic.

Prioritize correctness, clarity, maintainability, testability, and simplicity,
in that order. Avoid complexity that adds little value to this project.

## Project requirements

### Operations

- Addition, subtraction, multiplication, and division are required.
- Support negative numbers, zero, and decimal operands.
- Use two operand fields and an operation selector.
- Exponentiation, square root, and percentage are optional. Discuss their behavior
  only after all required work is implemented, tested, documented, and reviewed.

### Frontend

- Use React with TypeScript.
- Provide an intuitive interface for entering values and viewing results.
- Call the backend API for every calculation.
- Include input validation, error handling, and basic mobile responsiveness.
- Keep the code readable and cover important interactions with unit tests.

### Backend

- Use Go and expose REST endpoints for the four operations.
- Validate input and return results or errors as JSON.
- Handle invalid data, zero divisors, and relevant numeric edge cases.
- Keep arithmetic separate from HTTP behavior and server initialization.
- Cover important behavior with unit tests.

### Tests and coverage

Test happy paths, meaningful edge cases, and error handling in both applications.
Use simple assertions about observable behavior. Generate coverage reports without
chasing an artificial 100% target. Include relevant tests within each requirement.

### Docker

Provide multi-stage Dockerfiles for the backend and frontend, plus Docker Compose
so both start with one command. Serve the compiled frontend through Nginx and
proxy API requests to Go. Keep the configuration small.

### Documentation

Keep the README concise, visual, and approachable, starting with the Docker setup.
Include:

- Setup instructions and how to run both applications.
- Tests, builds, and coverage commands.
- API examples.
- Design decisions, assumptions, and a brief project layout.

Another developer should be able to run the project using only the README.
Use English for code, UI text, API messages, and project documentation.

### AI usage notes

`PROMPTS.md` retains the existing 21 entries. Do not record new prompts or append
new entries unless the user explicitly requests it. Keep project decisions and
development rules current in `AGENTS.md`.

## Planning and development agreement

The initial Grill Me discussion resolves requirement ambiguities, architecture,
API design, frontend/backend responsibilities, validation, errors, edge cases,
testing, Docker, and important technical choices before implementation.

Record approved decisions in the root `AGENTS.md`. It governs architecture,
conventions, commands, and the implementation process. Use GSD only when explicitly
requested. Do not delegate to subagents without explicit authorization.

## Incremental workflow

Implement one requirement at a time:

1. Analyze the scope, expected behavior, and affected areas.
2. Implement only that requirement and necessary prerequisites.
3. Review correctness, readability, and responsibility boundaries.
4. Add or adjust meaningful tests.
5. Run relevant tests and checks; fix failures and rerun affected checks.
6. Update relevant documentation. Record new prompts only when explicitly asked.
7. Report completion and stop for the user's review before the next requirement.

Do not anticipate future requirements without necessity or implement everything
in one pass.

## Engineering conventions

- Prefer clear names, small functions, low coupling, and explicit error handling.
- Remove meaningful duplication while keeping behavior easy to follow.
- Validate data at the appropriate boundary.
- Apply SOLID when it improves separation, maintenance, testing, or clarity.
- Do not introduce interfaces, factories, patterns, or layers just to demonstrate them.
- Do not invent product requirements or change unrelated code.
- Add dependencies only when there is a clear need.
- Keep frontend and backend responsibilities separate.
- Keep `AGENTS.md` current when an approved decision changes.
- Do not publish or perform GitHub writes without explicit authorization.

## Agreed implementation sequence

1. Minimal backend foundation.
2. Addition API with validation, errors, and tests.
3. Minimal frontend foundation.
4. Frontend addition integration, interaction states, and tests.
5. Subtraction, multiplication, and division, one at a time.
6. Interface, error handling, and responsiveness review.
7. Meaningful test review and coverage reports.
8. Docker images and Compose, developed incrementally.
9. Documentation and mandatory final review.
10. Discuss optional operations after required completion.

## Definition of done

Before advancing, confirm expected behavior, clear code, passing relevant tests,
considered edge cases, current documentation, and compliance with the approved
scope. Do not add unnecessary complexity. Give the user a completion checkpoint.

## Mandatory final review

Before declaring the project complete:

1. Run all frontend and backend tests, builds, and static checks.
2. Generate and inspect both coverage reports.
3. Review numeric, input, API, and communication edge cases.
4. Compare the implementation with `REQUISITES.md` and the approved decisions.
5. Identify and fix missing requirements and bugs.
6. Review Clean Code, useful SOLID application, and unnecessary complexity.
7. Build and run Docker from a clean state and verify integrated behavior.
8. Verify README setup instructions, commands, and API examples.
9. Confirm another developer can run the project using only the README.
10. Review existing AI notes for accuracy without adding new prompts automatically.
11. Report completion and remaining limitations before discussing optional work.
