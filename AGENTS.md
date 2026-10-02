# Project Development Agreement

## Purpose and sources of truth

Build a full-stack calculator using React, TypeScript, and Go. The frontend
must call the backend REST API to perform arithmetic operations.

Prioritize, in this order:

1. Correctness.
2. Clarity.
3. Maintainability.
4. Testability.
5. Simplicity.

`REQUISITES.md` defines the project requirements. `TASK.md` describes the
requested development process. This file records the decisions approved during
the Grill Me discussion and governs implementation. Preserve the original
project requirements; do not invent new product requirements.

Keep the required solution small and easy to explain. Go and Docker are part
of the approved project scope. Present the project as a generic calculator;
use generic project language in all Markdown files.

All code, identifiers, comments, UI text, API messages, and project documentation
must be in English. Discussion with the user may remain in Portuguese.

## Approved scope

- Implement addition, subtraction, multiplication, and division.
- Use two operand fields, an operation selector, and a Calculate button.
- Do not implement expression parsing, precedence, or chained operations.
- Support negative numbers, zero, and decimal operands.
- Provide input validation, error handling, and basic mobile responsiveness.
- Include relevant frontend and backend unit tests and coverage reports.
- Include Dockerfiles for both applications and Docker Compose.
- Provide a README with setup, execution, tests, coverage, API examples, design
  decisions, assumptions, and project structure.
- Preserve the existing 22 entries in `PROMPTS.md`. The prompt log is closed;
  add entries only if the user explicitly reopens it.
- Consider exponentiation, square root, and percentage only after every mandatory
  requirement is implemented, tested, documented, and reviewed. Define their
  behavior with the user before implementation.

## Architecture and dependencies

Use one repository with clearly separated `backend/` and `frontend/` directories.
No database or persistence is needed.

Backend:

- Use Go and its standard library, including `net/http`, `encoding/json`,
  `testing`, and `net/http/httptest`.
- Keep mathematical functions separate from HTTP decoding, validation, and
  response handling.
- Keep server initialization separate from mathematical and HTTP behavior.
- Use concrete functions and types. Introduce interfaces only when a demonstrated
  need improves clarity, maintenance, or testing.
- Go 1.26 is the initial toolchain baseline. Confirm and record the compatible
  patch and container image during setup.

Frontend:

- Use React, TypeScript, Vite, npm, plain CSS, and native `fetch`.
- Use React state for this single-screen application.
- Keep API access separate from UI rendering and interaction state.
- Extract components only when they improve readability or responsibility
  boundaries.
- Use Vitest and React Testing Library with a DOM test environment such as jsdom.
- Use an up-to-date Node 24 LTS patch compatible with the selected tooling. The
  installed Node 25 is not the approved runtime baseline.
- Pin compatible dependencies with `package-lock.json`; use `npm ci` for
  reproducible installations after the lockfile exists.
- Confirm exact package versions during setup. Do not assume that independently
  selected latest versions work together.

Do not add an HTTP framework, Axios, a router, global state management, a component
library, repositories, factories, or architectural patterns without a clear need.

### Planned structure

Create files only when the current requirement needs them. This is a target
structure, not an instruction to scaffold everything immediately.

```text
backend/
  go.mod
  main.go
  internal/
    calculator/          # Mathematical functions and related tests
    httpapi/             # HTTP handlers, validation, and related tests
  Dockerfile
frontend/
  package.json
  package-lock.json
  src/
    App.tsx
    api.ts
    ...                  # Styles and related tests as needed
  Dockerfile
  ...                    # Vite, TypeScript, test, and serving configuration
compose.yaml
README.md
PROMPTS.md
AGENTS.md
TASK.md
REQUISITES.md
```

## REST API contract

Expose these endpoints:

| Method | Path | Calculation |
| --- | --- | --- |
| POST | `/api/add` | left + right |
| POST | `/api/subtract` | left - right |
| POST | `/api/multiply` | left * right |
| POST | `/api/divide` | left / right |

All four endpoints receive a JSON object with numeric `left` and `right` fields.
Operand order matters for subtraction and division.

Example addition request:

```json
{"left":8,"right":2}
```

Successful response:

```json
{"result":10}
```

Errors use a consistent JSON envelope with stable English codes and clear English
messages. For example:

```json
{"error":{"code":"division_by_zero","message":"Cannot divide by zero."}}
```

| Status | Meaning |
| --- | --- |
| 200 | Successful calculation |
| 400 | Invalid input, division by zero, or a non-finite calculation result |
| 404 | Unknown path |
| 405 | Unsupported method; include the appropriate `Allow` header |
| 413 | Request body exceeds 1 KiB (1,024 bytes) |
| 415 | Missing or incompatible request content type |
| 500 | Unexpected internal failure |

Return JSON errors consistently, including routing and method errors. Do not expose
internal implementation details in user-facing errors. Do not manufacture extra
failure mechanisms solely to exercise a 500 response.

### Backend validation and numeric behavior

- Require `Content-Type: application/json`; valid media-type parameters may be
  accepted. Reject multiple Content-Type headers to avoid ambiguous input.
- Limit request bodies to 1,024 bytes.
- Require exactly one JSON object containing both operand fields.
- Require each operand field to occur once, including equivalent escaped names.
  Reject duplicate fields rather than allowing a later value to overwrite one.
- Reject malformed JSON, missing fields, unknown fields, `null`, strings,
  booleans, arrays, and other incorrect operand or body types.
- Reject any non-whitespace content after the JSON object.
- Distinguish a missing operand from a legitimate numeric zero.
- Accept only finite operands representable as Go `float64`.
- Reject division by positive or negative zero.
- Reject non-finite results, including overflow from finite operands.
- Do not add arbitrary minimum or maximum operand values beyond the finite
  `float64` representation.
- Return the computed value without additional business-level rounding.
- Document ordinary binary floating-point limitations. Exact decimal or financial
  arithmetic is outside the approved scope.

Backend validation is authoritative and must also protect direct API callers.

## Frontend responsibilities and behavior

The frontend collects input, validates it for usability, manages interaction
state, and displays API results and errors. All arithmetic runs on the backend.
Do not calculate locally as a fallback when the API fails.

### Input handling

- Accept negative numbers, decimals using a period, and scientific notation such
  as `1e3`.
- Ignore whitespace at the beginning and end of a field.
- Reject empty fields, commas, thousands separators, hexadecimal, non-finite
  values, and partially numeric text.
- Validate the entire input; do not silently accept a numeric prefix.
- Explain the period decimal separator in English.
- Distinguish invalid numeric format from a syntactically valid number outside
  the finite numeric range. Identify the first invalid field in the message.
  Preserve ordinary underflow to zero and do not add arbitrary operand limits.
- Send numeric operands in JSON after successful validation.
- Keep input parsing and result formatting separate from arithmetic.

### Results and request state

- Present results with at most 12 significant digits. Scientific notation may be
  used for very large or very small results.
- Preserve the backend value; presentation rounding must not alter API behavior.
- For example, the API may return `0.30000000000000004` while the UI displays `0.3`.
- Prefix results with `≈` and show an English precision note when the 12-digit
  presentation changes the backend value or an operand is an integer outside
  JavaScript's safe integer range. This is a conservative precision-risk signal;
  absence of the marker does not guarantee exact arithmetic. Preserve API values
  and associate the note with the accessible output description.
- Clear the previous result and error when an operand or the operation changes.
- Clear stale output when submitting a new calculation.
- Disable both fields, the operation selector, and the button during a request.
- Show `Calculating…` and allow only one calculation request at a time.
- Cancel waiting after 10 seconds and restore interaction after every outcome.
- Display known API error messages; use clear English messages for network
  failures, timeouts, and invalid API responses.
- Validate the expected response structure rather than trusting any JSON body.
- Allow manual retry after failure. Do not implement automatic retries.
- Use labeled controls, readable feedback, and a layout suitable for narrow mobile
  screens.
- Client validation focuses the first invalid operand and links the message to
  that control's accessible description. Clear the invalid mark along with stale
  feedback. After a request finishes, return focus to Calculate for keyboard retry.

## Testing and validation

Add or adjust relevant tests within each requirement. Do not defer its error
handling or tests to a later generic validation phase.

Backend tests must cover:

- Each operation's happy path and relevant negative, decimal, and zero cases.
- Division by positive and negative zero and non-finite calculation results.
- The HTTP contract: valid requests, malformed data, missing and incorrectly typed
  operands, unknown fields, trailing content, content types, body limits, methods,
  status codes, headers, and JSON response shapes.
- Floating-point assertions appropriate to the calculation; do not demand exact
  decimal equality where the representation does not support it.

Frontend tests must cover:

- Operand entry and operation selection.
- Client validation and the request sent to the API.
- Loading and disabled states, successful results, and presentation formatting.
- Clearing stale output when input changes.
- API errors, unavailable backend, timeout, malformed responses, and manual retry.

Use simple tests that verify observable behavior. Generate frontend and backend
coverage reports, but do not impose an artificial minimum or pursue 100% through
irrelevant tests. A browser end-to-end test suite is not part of the initial scope;
perform manual integrated and mobile-layout checks.

### Planned command interface

Backend run, build, test, static-check, and formatting support is implemented.
Frontend install, development server, tests, and build support is implemented in
requirement 3. Coverage report commands are implemented in requirement 7. The
backend Docker build/run commands and frontend/Compose setup are implemented
in requirement 8.
Create command support incrementally and verify it before documenting it as runnable.
Commands are relative to the indicated working directory.

| Directory | Purpose | Command |
| --- | --- | --- |
| `backend/` | Run the API | `go run .` |
| `backend/` | Build | `go build -o bin/calculator .` |
| `backend/` | Tests | `go test ./...` |
| `backend/` | Static checks | `go vet ./...` |
| `backend/` | Check formatting | `gofmt -l .` |
| `backend/` | Coverage data | `go test -coverprofile=coverage.out ./...` |
| `backend/` | Coverage summary | `go tool cover -func=coverage.out` |
| `backend/` | HTML coverage report | `go tool cover -html=coverage.out -o coverage.html` |
| `frontend/` | Install locked dependencies | `npm ci` |
| `frontend/` | Development server | `npm run dev` |
| `frontend/` | Type-check and production build | `npm run build` |
| `frontend/` | Run tests once | `npm run test -- --run` |
| `frontend/` | Coverage, including HTML | `npm run test:coverage` |
| Repository root | Build backend container image | `docker build -t fullstack-calculator-backend:local ./backend` |
| Repository root | Run backend container | `docker run --rm --name calculator-backend -p 127.0.0.1:18080:8080 fullstack-calculator-backend:local` |
| Repository root | Build and start both containers | `docker compose up --build` |
| Repository root | Stop both containers | `docker compose down` |

Keep generated binaries, coverage output, dependencies, and build output out of
version control. Record report locations and verified commands in the README when
implemented.

## Docker and frontend/backend communication

- Use two containers coordinated by Docker Compose.
- Use multi-stage builds for each application.
- Run the compiled Go backend in its runtime container.
- Serve the compiled frontend through Nginx.
- Have the frontend call relative `/api/...` paths.
- Configure Vite to proxy `/api` during development.
- Configure Nginx to proxy `/api` to the backend service in Docker.
- Preserve the agreed API paths through both proxies.
- This browser communication flow does not require CORS configuration.
- Keep addresses and ports in configuration rather than hardcoding a backend URL
  in React code. Document the selected ports during setup.
- Keep configuration small and ensure a clean build and start are reproducible
  from the README.

## Incremental implementation protocol

Use GSD only when the user explicitly requests it. Do not spawn subagents unless
the user or an applicable instruction explicitly authorizes delegation.

For each requirement:

1. Normalize the request into goal, constraints, relevant context, affected areas,
   implementation plan, and validation steps.
2. Analyze the exact scope and define expected behavior.
3. Implement only that requirement and its necessary prerequisites.
4. Review correctness, readability, and responsibility boundaries.
5. Create or adjust its relevant tests.
6. Run the relevant tests and checks.
7. Correct discovered problems and rerun affected checks.
8. Confirm its definition of done and update relevant documentation. Update the
   prompt log only when explicitly requested.
9. Report the outcome, validation evidence, and any limitations to the user.
10. Stop and wait for the user to authorize the next requirement.

Do not implement the entire project in one pass or anticipate future requirements
without necessity. Do not change unrelated code. Do not claim a test, build, or
Docker verification passed unless it was actually executed successfully.

### Agreed initial sequence

1. Minimal backend foundation.
2. Addition API, including validation, errors, and relevant tests.
3. Minimal frontend foundation.
4. Addition integration from frontend to backend, including relevant UI states and
   tests.
5. Remaining mandatory operations, one at a time, including related tests.
6. Complete outstanding interface, error handling, and responsiveness requirements.
7. Review remaining test coverage and generate coverage reports.
8. Implement Docker incrementally.
9. Complete documentation and perform the mandatory final review.
10. Discuss optional operations only after mandatory completion.

Keep documentation current throughout, rather than postponing it to the last
step. Record new prompts only when explicitly requested. Split a step into smaller requirements when
needed, and agree any material sequence changes with the user.

### Definition of done for each requirement

- Expected behavior works and matches the project requirements and approved decisions.
- Code is clear, idiomatic, and easy to explain.
- Relevant tests exist and pass.
- Relevant edge cases have been considered.
- No unnecessary complexity or unrelated changes were added.
- Related documentation is current. New AI prompt records require an explicit request.
- The user has received the completion checkpoint before the next requirement.

## Engineering conventions

- Use clear names, small functions, explicit error handling, predictable behavior,
  well-defined responsibilities, and low coupling.
- Apply SOLID only when it improves separation, testing, maintenance, or clarity.
- Avoid clever code, speculative interfaces, and abstractions added to demonstrate
  patterns.
- Remove meaningful duplication without sacrificing readability.
- Validate data at the appropriate boundary.
- Do not add dependencies without a clear need.
- Keep TypeScript checks and Go formatting part of the relevant validation.
- Update this file if an approved technical decision changes.
- Do not perform GitHub write operations without explicit user approval.
- Obtain explicit authorization before pushing or publishing the repository.

## AI prompt documentation

`PROMPTS.md` retains 22 entries and is closed after the final edge-case and
Clean Code review prompt. Do not append entries unless the user explicitly
reopens the log. User-requested edits to existing
entries are allowed, including generic project wording. The rules below apply
when prompt recording is explicitly requested:

- Translate Portuguese prompts into English; exact word-for-word fidelity is not
  required.
- Expand an approval such as "I accept the recommendation" into the complete
  technical decision the user approved.
- Include concrete behaviors, values, constraints, and error contracts, rather
  than recording generic approvals or references to question numbers alone.
- For example, record the full agreed 1 KiB body limit and HTTP status behavior
  instead of "Q21: recommendation accepted."
- Preserve the meaning and chronology of the actual approved work. Do not invent
  user instructions or misrepresent suggestions that were not approved.
- Omit irrelevant conversational acknowledgments while keeping substantive AI
  usage transparent.

## Mandatory final review

Before declaring the mandatory project complete:

1. Run all frontend and backend tests.
2. Run the relevant builds and static checks.
3. Generate and inspect both coverage reports.
4. Review numeric, input, API, and communication edge cases.
5. Compare the result with `REQUISITES.md` and `TASK.md` requirement by requirement.
6. Identify and correct missing requirements and bugs.
7. Review Clean Code and useful SOLID application; remove unnecessary complexity.
8. Build and run Docker from a clean state and verify integrated behavior.
9. Verify the README's setup, execution, test, coverage, and API examples.
10. Check that another developer can clone and run the project using only the
    README's documented prerequisites and instructions.
11. Review the existing `PROMPTS.md` for accuracy; do not add new prompts automatically.
12. Report completion and remaining limitations before discussing optional work.

## Current checkpoint

Requirements 1 and 2 are complete: the minimal backend foundation and the addition
API. The backend defaults to port 8080, configurable through `PORT`.

`POST /api/add` implements the approved request, response, validation, and error
contracts. Addition lives in `internal/calculator`; HTTP validation and responses
live in `internal/httpapi`. At that checkpoint, no other arithmetic operation had
been implemented.

Backend tests, build, `go vet`, and formatting checks passed with Go 1.26.0. Eight
local HTTP checks verified addition, the exact size limit, invalid input,
overflow, incorrect method, unsupported content type, excessive size, and unknown
path on a temporary server, which was stopped after validation.

A subsequent test-quality review strengthened the decimal assertion to reject
non-finite results and removed three redundant HTTP success cases. Tests detected
all 14 deliberately introduced, selected regressions in a temporary copy through
assertion failures. Actual project tests passed again with `go test -count=1 ./...`,
along with static analysis and formatting checks. This audit did not change
production behavior or authorize the next requirement.

Requirement 3, the minimal frontend foundation, is also complete. It includes the
React and HTML entry points, strict TypeScript configuration, plain CSS application
shell, development server, locked dependencies, and Vitest/Testing Library/jsdom
setup. Node 24.19.0 was used to verify clean installation from the lockfile, the
single rendering smoke test, type checks, and production build. The screen was
inspected in the browser in both development and production preview modes. The
temporary production preview server was stopped. The development frontend uses
port 5173 with strict port
selection; the backend remains on port 8080 by default.

The frontend `.nvmrc` pins Node 24.19.0; package engines require a Node 24 patch at
least that new, and `.npmrc` enforces it during dependency installation. The
assistant used an existing bundled runtime without changing the user's global
Node installation. Frontend commands in a user terminal also require Node 24.

Requirement 4, frontend-to-backend addition integration, is complete. The form
validates full numeric input and calls the separate API client through relative
`/api/add`. Vite proxies `/api` to the backend on port 8080 by default; `BACKEND_URL`
in the frontend environment can override the development proxy target. The UI
handles pending requests, a ten-second timeout, response validation, errors,
manual retry, stale output, zero results, and presentation with at most 12
significant digits. Only Addition is available in the selector at this checkpoint.

All 54 frontend tests, strict TypeScript checks, and the production build passed
on Node 24.19.0. Browser checks against the real Go backend through Vite verified
integer and decimal addition, zero, clearing output, invalid input, overflow,
and retry. The form was inspected at a 320-pixel viewport width. No backend code
changed in this requirement. README and PROMPTS.md are current.

The first operation in requirement 5, subtraction, is also complete.
`POST /api/subtract` computes `left - right` and shares HTTP validation and
response handling with addition. Mathematical functions remain separate; their
finite-number check and error values live in `internal/calculator/numbers.go`.
The frontend has a controlled Addition/Subtraction selector and a shared API
client accepting only the implemented operations. Changing the operation clears
both results and errors. At that checkpoint, multiplication and division were
not yet implemented.

Backend tests, build, static checks, and formatting passed. Frontend tests
(57 cases), strict TypeScript checks, and production build passed on Node 24.19.0.
Tests cover subtraction order, negative and zero operands, decimals, finite
boundaries, overflow in both directions, and invalid operands. Shared HTTP
contract tests exercise both endpoints. Browser checks used temporary backend
and frontend processes on ports 18080 and 5174 to verify both operand orders,
decimal presentation, overflow, operation changes, and addition regression
behavior. These temporary processes were stopped afterward; the user's running
servers were left alone. An existing Go process must be restarted to load the
new endpoint. README and PROMPTS.md are current.

Multiplication in requirement 5 is now complete. `POST /api/multiply` uses the
separate `calculator.Multiply` function and the existing HTTP processing.
Multiplication is available in the frontend selector and shared API client.
Finite operands and results are required; ordinary float64 underflow to zero is
accepted, consistently with the approved numeric contract.

Backend tests, build, static checks, and formatting passed. All 59 frontend tests,
strict TypeScript checks, and the production build passed on Node 24.19.0.
Mathematical tests cover signs, zero on both sides, decimals, finite boundaries,
underflow, invalid operands, and overflow in both directions. Shared HTTP
validation tests now cover all three endpoints. New frontend tests verify the
multiplication endpoint, numeric payload, selection, result, and clearing output.
Existing request-state and formatting tests remain shared.

Browser checks with temporary servers on ports 18080 and 5174 verified positive
and negative multiplication, zero, decimal presentation, overflow, and manual
retry. Both temporary servers were stopped afterward. The user's servers were
left alone; an existing Go process must be restarted to load the new endpoint.
README and PROMPTS.md are current.

The user's follow-up input-validation refinement is complete. The parser returns
either a finite numeric value or an explicit invalid-format/out-of-range reason.
The UI identifies First number or Second number in the message, rejects the input
before fetch, and clears feedback when an operand changes. This changes feedback
without changing the accepted numeric range or backend behavior. All 68 frontend
tests, strict TypeScript checks, and the build passed. Browser checks confirmed
format and range messages and second-field identification. Documentation is current.

Division is complete, concluding requirement 5's remaining mandatory operations.
`POST /api/divide` computes `left / right` in a separate mathematical function.
The backend rejects positive and negative zero divisors, including zero divided
by zero, with HTTP 400 and `division_by_zero` / `Cannot divide by zero.`
Finite-result validation and the existing HTTP contract remain shared. The
frontend exposes Division and displays backend errors with manual retry; no
arithmetic is performed locally.

Backend tests, build, static checks, and formatting passed. All 72 frontend
tests, strict TypeScript checks, and the production build passed on Node 24.19.0.
Tests cover operand order, signs, zero numerators, positive/negative zero
divisors, decimals, very small nonzero divisors, underflow, invalid operands,
and overflow in both directions. Shared HTTP tests now cover all four routes;
success-response assertions were consolidated without removing behavior checks.
Frontend tests verify selection, operand order, both zero-input spellings,
backend feedback, control restoration, and retry.

Browser checks on temporary servers verified both operand orders, one-third
presentation, zero and negative-zero input, overflow, and retry. A direct HTTP
check confirmed raw JSON `-0` returns the approved 400 error envelope. Temporary
servers on ports 18080 and 5174 were stopped; the user's servers were left alone.
Restart an existing Go process to load division. Documentation is current.

Requirement 6, the interface, error handling, and responsiveness review, is
complete. Review found two keyboard/feedback gaps: client validation left focus
on Calculate and did not associate the error with its field; disabling controls
during requests lost keyboard focus. Client errors now focus the first invalid
operand, mark it with aria-invalid and a red border, and link the message through
aria-describedby. API errors remain general feedback. Request completion restores
focus to Calculate. No arithmetic, API contract, or dependency changed.

All 73 frontend tests, strict TypeScript checks, and the production build passed
on Node 24.19.0. Tests were strengthened for focus and accessible field feedback,
clearing the invalid mark, and focus restoration after success, API failure,
and timeout. The existing loading, network, malformed-response, formatting,
validation, and retry checks also passed. No backend code changed in this requirement.

Manual browser checks verified Tab order, Enter submission, focus on invalid
operands, division-by-zero feedback, and successful retry with the real backend.
Viewport checks at 320x740, 375x812, 768x360, and 1280x800 confirmed field stacking
or side-by-side layout, readable long errors, large scientific-notation results,
and no horizontal overflow. The temporary viewport override was reset. These
checks do not claim physical-device or screen-reader speech verification.
README and PROMPTS.md are current.

Requirement 7, meaningful test review and coverage reports, is complete.
The frontend uses @vitest/coverage-v8 5.0.3, matching Vitest, with the verified
npm run test:coverage command. The report includes application TypeScript,
including the unimported main.tsx; only tests and test setup are excluded.
Go coverage includes all backend packages. No coverage threshold was imposed.

Review added a real-handler request-body read failure test and separated invalid
API error-envelope cases so incorrect types and blank code/message fields are
checked independently. A JSON 503 response without the error envelope also
verifies unavailable-service feedback. Production calculation and UI behavior
were not changed. All 77 frontend tests, strict TypeScript checks, and production
build passed on Node 24.19.0. Backend coverage tests, go vet, and formatting passed.
An isolated clean npm ci installation verified the updated dependency lockfile.

Generated reports: backend/coverage.out and backend/coverage.html;
frontend/coverage/index.html, coverage-final.json, and coverage-summary.json.
All are excluded from version control. Backend statement coverage is 90.6%
overall, 100% for calculator, and 96.2% for httpapi. Frontend statement coverage
is 95.34% and branch coverage is 94.82%; api.ts and numbers.ts have 100% on both.
App.tsx has 100% statements and 95.83% branches. main.go and main.tsx retain 0%
unit coverage; startup was exercised in prior integrated manual checks.
Remaining defensive branches are documented in the README. Both HTML reports
were opened and inspected in the browser.

Six selected regressions in temporary copies were detected by assertion
failures: incorrect subtraction/multiplication, reversed division operands,
missing zero-divisor rejection, and acceptance of blank API error code/message.
This was a bounded test-quality check, not an exhaustive mutation score. Temporary
copies and the report-inspection server were removed/stopped. README and
PROMPTS.md are current.

The first requirement 8 checkpoint, the backend Docker image, is complete.
backend/Dockerfile uses golang:1.26.8-alpine3.24 in its build stage, a verified
compatible patch of the Go 1.26 baseline. All backend tests run before compiling
with CGO_ENABLED=0. The scratch runtime contains the executable only, runs as
UID/GID 65532:65532, and starts it directly. PORT defaults to 8080 and remains
configurable. backend/.dockerignore excludes local binaries, coverage files,
Git metadata, and environment files. No application code changed.

A docker build --pull --no-cache passed on Docker Desktop 4.44.3 / Engine 28.3.2,
including all tests with the container's Go 1.26.8 toolchain. Fifteen HTTP checks
against temporary containers passed: all four operations, invalid input,
positive/negative zero divisors, overflow, unknown path, incorrect method and
Allow, missing/incompatible content type, exact/excessive request-body size, and
PORT=8081. Checks verified JSON response shapes and content types. Inspection
confirmed the configured and running UID 65532 and the single runtime filesystem
layer. Local Go static analysis and formatting passed. The validated platform is
Linux ARM64; AMD64 has not been verified. Temporary containers were stopped and
removed, leaving the built image for reuse. Existing development servers were
not changed. README and PROMPTS.md are current.

The documented standalone mapping is localhost:18080 to container port 8080 so
it can coexist with the development API. Frontend code continues to use relative
API paths; the existing BACKEND_URL setting can point Vite at this container.

Requirement 8 is complete after the user explicitly authorized implementing the
frontend Docker image and Compose together to provide one startup command.
frontend/Dockerfile uses node:24.19.0-alpine3.24 for npm ci, all frontend tests,
and the TypeScript/Vite build. nginx:1.30.5-alpine3.24 serves the compiled assets.
The frontend build context excludes node_modules, dist, coverage, Git metadata,
and environment files. No application source or dependency changed.

compose.yaml builds both services and publishes only frontend port 80 at
127.0.0.1:8081 by default; APP_PORT can change that host port. Backend PORT is
8080 on the internal Compose network and has no host port mapping. The single
startup command from the root is docker compose up --build. docker compose down
stops and removes the application's containers and network.

frontend/nginx.conf preserves API paths, backend errors, and Allow headers.
The Go handler owns the 1024-byte limit, with Nginx API body-size limiting disabled
and request buffering disabled so excessive bodies retain the JSON 413 contract.
Proxy retries are disabled. A resolving upstream uses Docker DNS at 127.0.0.11
with a five-second validity to follow backend address changes. Compose controls
start order, not application readiness; no new health endpoint or wait script
was added. Proxy failures use the existing frontend unavailable-service handling
and allow manual retry.

Both images built with --pull --no-cache. All backend tests and all 77 frontend
tests passed inside their builds, including strict TypeScript checks and the
production frontend build. Compose configuration and Nginx syntax checks passed.
Compiled HTML/JS/CSS loaded, and 17 API checks through Nginx passed, including all
operations, numeric/input errors, methods, paths, content types, Allow, the exact
1024-byte boundary, 1025 bytes, and a body above 1 MiB. Inspection confirmed the
backend had no published port. Browser checks verified all operations, decimal
presentation, division-by-zero feedback and corrected-input retry, unavailable
backend feedback, and successful manual retry after backend recreation without
recreating the frontend. Validation used Linux ARM64; AMD64 is not verified.
An APP_PORT=8082 startup check also passed. Docker Compose down removed the
validation containers and network; images remain available for reuse. Existing
native development servers were left running. README and PROMPTS.md are current.

The user-approved approximation feedback refinement is complete. Result state
keeps the raw API value plus a feedback flag. needsApproximationNotice compares
presentation with that value and flags operands that are integers outside the
JavaScript safe integer range. Approximate results show an approximately-equal
symbol and the English limited-precision note, linked to the output through
aria-describedby. The note clears with stale output. Ordinary results such as
8 + 2 show no marker. Large-integer detection is conservative and does not claim
exhaustive detection of all binary rounding or underflow. No backend code,
numeric range, API response, or dependency changed.

All 88 frontend tests, coverage, strict TypeScript checks, and production build
passed on Node 24.19.0. The updated frontend image also ran all 88 tests and built
successfully. Helper/UI tests cover presentation rounding, repeating fractions,
safe and unsafe integer boundaries, both operand positions and signs, accessible
precision feedback, unchanged payloads, and clearing on edits, operation changes,
new requests, and failure. Current frontend coverage is 95.6% statements and
95.58% branches; App has 100% statements and 96.66% branches, while API/numeric
helpers retain 100% on both. main.tsx remains included with zero unit coverage.

An isolated Compose application on port 8082 verified 8 + 2 without the marker,
0.1 + 0.2 with approximate feedback, and the user's large-integer example. A
320x740 viewport check found no horizontal overflow, and the override was reset.
Temporary validation containers were removed. The user's Compose application on
port 8081 was not restarted; it needs a rebuild and browser refresh for the new
compiled frontend. README and PROMPTS.md are current.

The requested Everyday math subtitle removal is complete; its unused CSS rule
was removed too. All 24 form tests, TypeScript checks, and the production build
passed. The running Docker frontend needs rebuilding for this source change.

The README is now a concise Docker-first guide with native setup, API examples,
tests, coverage, design decisions, and a brief project layout. TASK.md and
REQUISITES.md use generic English project language. Automatic prompt recording
is disabled; the existing 21 entries remain, with only requested wording edits.
The mandatory final review is still pending.

The requested input and numeric edge-case audit is complete. Regression tests
first reproduced two validation gaps on all four endpoints: duplicate JSON
operand fields and ambiguous multiple Content-Type headers. Shared HTTP validation
now rejects them with 400 / invalid_input and 415 / unsupported_media_type,
respectively. JSON decoding checks each field before accepting its value; escaped
duplicate names and null values overwritten by numbers cannot bypass validation.

Expanded tests cover malformed JSON numeric syntax, incomplete exponents, comments,
trailing non-JSON whitespace, HEAD/Allow behavior, finite float64 boundaries,
tiny nonzero divisors, overflow, and positive/negative underflow to zero. Frontend
tests cover empty fields, isolated signs, expressions, internal whitespace,
invisible characters, unsupported Unicode numeric symbols, signed zero, and
extreme-result formatting. No numeric range, arithmetic, or UI behavior changed.

All 120 frontend tests, strict TypeScript checks, and production build passed.
Backend tests, build, go vet, and formatting passed. Updated coverage is 92.2%
overall backend statements and 97.1% for HTTP handlers; arithmetic remains 100%.
Frontend coverage is 95.6% statements and 95.58% branches. Generated reports
remain excluded from Git. The new Go fuzz target completed 618,050 executions
in a bounded 20-second run without failures; this does not claim exhaustive
input or transport testing.

An isolated Compose application on port 18081 rebuilt both images, including their
tests, and passed 256 HTTP checks through Nginx. Checks included all operations,
invalid inputs, duplicate fields, numeric extremes, exact/excessive body sizes,
chunked oversized requests, media types, methods, paths, and 24 concurrent
requests. The temporary containers and network were removed. The user's existing
application was not restarted; it needs rebuilding to load backend changes.
PROMPTS.md was not changed. This focused audit does not complete the entire
mandatory final project review.

Requirement 9, the mandatory final project review, is complete following the
user's request to audit test quality and compare the application with the supplied
requirements. All mandatory features and deliverables are implemented locally:
the four operations, React/TypeScript frontend, Go REST API, validation, JSON
results/errors, basic mobile responsiveness, tests, generated coverage reports,
and concise setup/API/design documentation. Both Dockerfiles and unified Compose
startup also fulfill the optional Docker requirement. Advanced operations remain
optional and require a separate behavior discussion and authorization.

Review removed one duplicate HTTP tiny-divisor success case; the same request and
assertion remain covered in the division endpoint tests. The API-client test now
returns a deliberately different mock result to detect local recalculation as
well as unwanted rounding. Eight selected regressions in temporary copies were
detected by assertion failures: incorrect addition, missing zero-divisor rejection,
acceptance of duplicate fields or ambiguous media types, local frontend arithmetic,
numeric-prefix acceptance, enabled pending controls, and retained stale results.
This is a bounded test-quality check, not an exhaustive mutation score.

Backend tests, build, go vet, and formatting passed. All 120 frontend tests,
TypeScript checks, and production build passed. Coverage reports were generated
and inspected: backend statements 92.2% overall, arithmetic 100%, HTTP handlers
97.1%; frontend statements 95.6% and branches 95.58%. Startup files remain included
with zero unit coverage and were exercised through actual native and Docker startup.

A clean copy of the 48 deliverable project files excluded dependencies, generated
artifacts, and local environment files. npm ci succeeded with Node 24.19.0, and
both Docker images built without build-cache reuse, running their tests. Compose
startup passed 256 HTTP checks, including 24 concurrent requests. Native Go/Vite
startup passed 11 additional API/proxy/README checks on temporary ports 18080/5174.
Browser checks verified four operations, keyboard submission, invalid-input focus,
division-by-zero retry, approximation feedback, unavailable-service feedback,
and successful manual retry after restarting the backend. Basic layout checks
at 320 and 375 pixels found no horizontal overflow. No physical-device or
screen-reader speech verification is claimed. Validation used Linux ARM64 Docker;
AMD64 is not verified. Ordinary float64 limitations remain documented.

Temporary servers, containers, network, browser tab, viewport override, and clean
copy were removed or reset. The user's running application was left alone.
PROMPTS.md remains closed at 22 entries. The user handles the push after the
authorized local commit of the validated changes.
Do not implement optional operations or publish without further authorization.

The post-review frontend test correction waits for request-completion focus
restoration in success and API-error assertions. Rendering a result or error does
not guarantee that React's focus effect has finished. Production code is unchanged.
All 120 tests passed in five consecutive coverage runs, and the TypeScript/Vite
build passed. Removing focus restoration in an isolated copy still failed the
corrected success test on its focus assertion. Coverage remains 95.6% statements
and 95.58% branches. PROMPTS.md remains closed.
