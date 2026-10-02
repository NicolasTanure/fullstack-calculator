# AI Development Prompts

The prompts below are recorded in English. Portuguese requests are translated,
and approvals of recommendations are expanded into the actual approved technical
decisions. They preserve the substance of the discussion rather than its exact
wording. Implementation notes are identified separately from user prompts.

## 1. Calculator planning request

Read the entire `TASK.md` at the repository root before implementing code. It
contains the calculator project brief, requirements, and the development process I
want to follow.

Conduct a Grill Me discussion covering requirement ambiguities, architecture, API
design, frontend/backend responsibilities, validation, edge cases, testing,
Docker, and important technical decisions. Do not implement anything yet.

After the discussion, use the approved decisions to create the project's
`AGENTS.md`. Then implement one requirement at a time: analyze, implement, review,
test, correct, finalize, and only then consider the next requirement.

## 2. Original requirements and scope decisions

Read the newly created `REQUISITES.md`, which contains the initial project
requirements. If chained operations are not required, use two numeric fields, an
operation selector, and a Calculate button. Do not add expression parsing or
chained calculations.

Use Go `float64` arithmetic and accept ordinary floating-point limitations rather
than introducing exact decimal arithmetic. Expose separate POST endpoints for
addition, subtraction, multiplication, and division. The backend is authoritative
for arithmetic and validation; the frontend owns input assistance and presentation.

Include validation, relevant errors, and tests within each requirement. Implement
the minimal backend first, then addition in the API, then the frontend foundation
and addition integration. Do not postpone correctness to a later validation phase.

Use English for code, UI, API messages, README, `AGENTS.md`, and `PROMPTS.md`.
Complete, test, document, and review mandatory requirements before considering
exponentiation, square root, or percentage. Optional operations remain desirable
if time permits after mandatory completion.

## 3. Architecture, tooling, and API decisions

Use one repository with `backend/` and `frontend/`, without a database or
persistence. Separate mathematical functions, HTTP handling, and server
initialization. Separate API access from React rendering and extract components
only when useful. Prefer concrete functions and a small number of layers.

Use the Go standard library for HTTP, JSON, and backend tests. Use React,
TypeScript, Vite, npm, plain CSS, and native `fetch`. Use Vitest and React Testing
Library with a DOM test environment. Use compatible locked dependencies and Node
24 LTS. Do not introduce frameworks or libraries without a clear need.

Implement `POST /api/add`, `POST /api/subtract`, `POST /api/multiply`, and
`POST /api/divide`. Requests use numeric `left` and `right` fields. Successful
responses use a numeric `result` field. Errors use an `error` object containing a
stable `code` and an English `message`.

Reject missing fields, null operands, numeric strings, booleans, unknown fields,
malformed JSON, and extra content after the JSON object. Accept only finite
numeric operands, allow negatives and zero, reject division by positive or
negative zero, and report non-finite results explicitly. Do not impose arbitrary
operand bounds.

Return unrounded floating-point results from the backend. The UI should display
at most 12 significant digits and may use scientific notation. For example, an
API result of `0.30000000000000004` can appear as `0.3` in the interface.

Use two containers with multi-stage builds and Docker Compose. Serve the frontend
through Nginx and proxy relative `/api` requests to the backend. Use a Vite proxy
during development, avoiding a hardcoded backend URL in React.

Test mathematical rules, the HTTP contract, and important frontend interactions,
including happy paths, invalid inputs, loading, and failures. Generate coverage
reports without an artificial minimum. Manually verify the integrated application
and Docker; do not add a browser end-to-end suite in the initial delivery.

## 4. Input, interaction, and HTTP decisions

Accept negative numbers, period decimals, and scientific notation such as `1e3` in
the frontend. Trim leading and trailing whitespace. Reject empty input, commas,
thousands separators, hexadecimal, non-finite numbers, and partially numeric
text. Explain the period decimal separator to the user.

While calculating, disable both fields, the operation selector, and the button,
and show `Calculating…`. Clear previous output when the inputs or operation change.
Release the controls after success or failure. Cancel waiting after 10 seconds,
show clear English errors for network failures or invalid responses, and permit
manual retry without automatic retries.

Use `200` for successful calculations and `400` for invalid inputs, division by
zero, or non-finite results. Use `413` for an oversized body, `415` for a missing or
incompatible content type, `404` for an unknown path, and `405` with the appropriate
`Allow` header for an unsupported method. Use `500` for unexpected internal
failures. Keep errors in the agreed JSON format. Require `application/json` and
limit request bodies to 1 KiB (1,024 bytes), which is enough for two operands.

Complete and verify one requirement at a time, then stop for my review before
starting the next. First present the consolidated planning agreement; after
confirmation, create only `AGENTS.md` without starting application code.

## 5. Prompt recording instructions

Translate substantive prompts into English; exact word-for-word fidelity is not
necessary. Do not record generic answers such as "I accept the recommendation."
Instead, expand each approval into the actual full technical decision, including
behaviors, limits, and error statuses. Keep the record detailed and do not invent
instructions or approvals that did not occur.

## 6. Planning agreement confirmation

Confirm the consolidated development agreement, including Node 24 LTS. Create
`AGENTS.md` with the approved decisions and planned commands. Do not begin the
first implementation requirement in that step.

## 7. First implementation requirement

Start the first requirement: the minimal backend foundation. Follow `AGENTS.md`,
analyze the scope, implement it, review it, add and run relevant tests, correct any
problems, and stop for my review before moving on to addition.

### AI implementation scope

The assistant defined this requirement as a Go module, configurable HTTP server,
and JSON 404 responses for unknown paths, with handler tests and initial
documentation. Port 8080 is the default, overridden by `PORT`. No arithmetic
endpoint, frontend, external dependency, or Docker scaffolding belongs to this
requirement.

### AI review and validation outcome

The assistant reviewed server initialization and the independently testable HTTP
handler. Three unknown-path cases verify 404 status, JSON content type, and the
complete error envelope. Backend tests, build, static analysis, and formatting
checks passed with Go 1.26.0. Local HTTP checks verified the compiled server on the
default and a configured port, and the temporary servers were stopped. Requirement
1 is complete; addition remains pending user authorization.

## 8. Addition requirement authorization

All tests for the minimal backend have passed when I ran them locally. Proceed
with the next requirement: the addition API. Follow the agreed incremental
process and API contract, including input validation, error handling, and relevant
tests. Complete and review this requirement before beginning the next one.

### AI implementation scope

Implement only `POST /api/add`, with numeric `left` and `right` operands and a
numeric `result` response. Put mathematical logic in `internal/calculator`,
separate from HTTP handling. Keep JSON errors consistent and preserve unknown-path
behavior. Enforce required fields, correct numeric types, exact field names,
finite operands and results, a single JSON object, and the 1,024-byte body limit.
Support valid JSON media-type parameters. Test normal calculations, zero,
negatives, decimals, floating-point boundary values, invalid data, overflow,
methods, content types, and both sides of the body-size boundary. Update the
README and development checkpoint without implementing the frontend or another
operation.

### AI review and validation outcome

Backend tests, build, static analysis, and formatting checks passed. The mathematical
tests include direct non-finite inputs and both positive and negative overflow.
HTTP tests cover successful and invalid requests and retain the original 404
tests. A compiled server on a temporary port passed eight local HTTP checks:
addition, exact body limit, invalid operand, overflow, wrong method, wrong content
type, excessive body size, and unknown path. The temporary server was stopped;
the user's existing development server was not modified.

The assistant documented that a running `go run .` process must be restarted to
load code changes. Requirement 2 is complete; the frontend foundation remains
pending user authorization.

## 9. Test-quality review request

The tests and addition API call passed when I ran them locally. Make sure our
tests provide real value and are not being written just for the sake of having
tests. Review what they actually assert, identify redundant cases or weak
assertions, and verify that they detect meaningful regressions before we advance.

### AI review and validation outcome

The assistant found that the approximate decimal assertion could incorrectly
accept NaN. An isolated experiment reproduced that false positive; the assertion
was changed to reject non-finite results explicitly. The complete original HTTP
suite already rejected the invalid response, so this was a weakness in the
individual mathematical assertion, not an unnoticed arithmetic production bug.

Three redundant HTTP success cases were removed. HTTP negative input, decimal
serialization, and zero operands remain covered; the full arithmetic cases remain
in the mathematical tests. Content-type test names were made descriptive.

In a temporary copy, the updated tests detected all 14 selected regressions through
assertion failures: subtraction instead of addition, NaN in the decimal case,
removed finite-operand validation, removed overflow validation, rejected zero
operands, accepted unknown fields, null treated as zero, ignored trailing JSON,
body limits of 1,025 and 1,023 bytes, accepted text/plain, removed method validation,
omitted Allow header, and truncated HTTP results. Compilation failures were not
counted as detection. The temporary copy was deleted after the audit.

The actual project tests were rerun without test-result caching using
`go test -count=1 ./...`; they passed, along with static analysis and formatting
checks. This was a focused audit without a permanent mutation-testing framework.
No production behavior or requirement scope changed, and the frontend remains
pending user authorization.

## 10. Frontend foundation authorization

Proceed with the next requirement after the backend test-quality review: the
minimal frontend foundation. Follow the approved React, TypeScript, Vite, npm,
and Node 24 decisions. Complete, review, and verify this requirement, then stop
before implementing the addition integration.

### AI implementation scope

Create only the frontend package, HTML and React entry points, strict TypeScript
configuration, an initial application shell with plain CSS, and the Vitest test
setup with a DOM environment. Do not introduce the calculator form, API client,
arithmetic, or integration in this step. Use one rendering smoke test to validate
the React and DOM test setup without inventing behavior just to increase test
counts. Confirm dependency engines and peer requirements, pin compatible versions,
generate the npm lockfile, verify installation from it, run tests and the production
build, and inspect the development screen in the browser.

### AI review and validation outcome

The assistant found an existing Node 24.19.0 runtime in the Codex environment and
used it for verification without replacing the user's globally installed Node 25.
Package versions, engines, and peer requirements were checked in the official npm
registry. The frontend pins React 19.3.0, TypeScript 7.0.2, Vite 8.3.2, and Vitest
5.0.3, with compatible test and type dependencies. `.nvmrc` records Node 24.19.0
and `.npmrc` enforces the supported Node range during dependency installation.

Installation succeeded, followed by a successful clean `npm ci` using the cached
packages. The one rendering smoke test passed, the strict TypeScript and production
build passed, and `npm ls --depth=0` reported the expected dependencies. The initial
screen rendered in the browser through the development server on port 5173. The
production build also rendered through `npm run preview` on port 4173; that test
server was stopped afterward. The development preview was retained for user
review. No backend code changed and no API
requests or calculator interactions were implemented. Documentation was updated;
requirement 4 remains pending user authorization.

## 11. Addition integration authorization

The frontend tests and production build passed. Confirm how to run and inspect
the development screen in the browser as part of reviewing the requirement.
After reviewing it, proceed with frontend-to-backend addition integration,
following the approved input, request-state, error-handling, and testing decisions.
Complete this requirement and stop before implementing another operation.

### AI implementation scope

Add two labeled operand fields, an Addition selector, and a Calculate button.
Parse complete finite numeric inputs, trim surrounding whitespace, and reject
invalid formats without sending a request. Keep numeric parsing and presentation
separate from the API client and React interaction state. Send numeric operands
to `/api/add`; configure Vite to proxy `/api` to the backend, with a configurable
`BACKEND_URL` target and an environment example. Keep arithmetic on the backend.

Clear stale output on input changes and submission. Disable all controls during
the request, enforce the ten-second deadline through response-body reading,
validate success and error response shapes, and restore interaction after every
outcome. Support manual retry and display zero results correctly. Format results
with at most 12 significant digits without changing the backend calculation.

### AI review and validation outcome

All 54 tests across numeric parsing/formatting, the API client, and React behavior
passed with Node 24.19.0. The strict TypeScript checks and production build also
passed. UI tests use the real API client while replacing only fetch; assertions
check observable inputs, requests, results, errors, disabled controls, cleared
output, and retry behavior. Timer tests cover both waiting for headers and reading
a stalled response body. No backend code changed.

Browser verification through the Vite proxy and running Go server confirmed
`8 + 2 = 10`, presentation of `0.1 + 0.2` as `0.3`, and `-8 + 8 = 0`.
Changing inputs cleared the result; a comma decimal was rejected; overflowing
addition displayed the backend error; correcting inputs allowed successful retry.
The form was visually inspected at a 320-pixel viewport width and the temporary
viewport override was reset afterward. The development preview remains available
for user review. Documentation and the development checkpoint were updated;
subtraction remains pending authorization.

## 12. Subtraction authorization

The addition integration looks good after review. Proceed with the next
requirement: subtraction through the backend and frontend. Follow the incremental
implementation, review, test, correction, and documentation process, and stop
before implementing multiplication.

### AI implementation scope

Implement `POST /api/subtract` as `left - right`, preserving operand order and
the approved JSON, validation, body-limit, method, and error contracts. Add
Subtraction to the UI selector. Share HTTP processing and frontend API request
handling now that two operations need the same behavior, while keeping each
mathematical function explicit and independently testable. Keep the existing
timeout, loading, validation, error, and presentation behavior. Clear stale
results and errors when the selected operation changes. Add no dependencies.

### AI review and validation outcome

Backend tests, build, static analysis, and formatting checks passed. The tests
verify reversed operands, negative numbers, zero, decimals, cancellation and
smallest finite differences, non-finite operands, and both overflow directions.
The same HTTP validation cases run against addition and subtraction, avoiding
copied validation suites while checking that both routes enforce the contract.

All 57 frontend tests passed, along with strict TypeScript checks and the
production build on Node 24.19.0. New assertions verify the subtraction endpoint,
unchanged operand order, operation selection, and clearing both results and errors
on operation changes. Existing addition and request-state tests also passed.

Browser checks against temporary backend/frontend servers verified `8 - 2 = 6`,
`2 - 8 = -6`, and presentation of `0.3 - 0.2` as `0.1`. Overflow showed the
backend error; switching to addition cleared it, and the same large opposite-sign
operands then produced zero. The temporary servers were stopped after verification.
The user's existing servers were not modified; an existing Go process needs a
restart to load the new endpoint. Documentation was updated. Multiplication
remains pending user authorization.

## 13. Multiplication authorization

The subtraction implementation looks good after review. Proceed with the next
requirement: multiplication in the backend API and frontend selector. Follow the
incremental development process and stop for review before implementing division.

### AI implementation scope

Add `POST /api/multiply` with the approved numeric input and JSON response/error
contracts. Keep multiplication in its own mathematical function and reuse the
existing HTTP validation and API client. Add only Multiplication to the operation
selector, without new dependencies or unrelated UI changes. Verify signs, zero,
decimals, finite limits, underflow, invalid operands, and overflow. Frontend tests
should verify the new selection and request path while retaining the existing
shared formatting, validation, timeout, loading, and error tests.

### AI review and validation outcome

Backend tests, build, static analysis, and formatting passed. The arithmetic
tests cover positive/negative operands, zero on either side, exact and inexact
decimals, maximum and minimum finite results, underflow to zero, non-finite
operands, and positive/negative overflow. The HTTP tests verify multiplication
results, unrounded serialization, and overflow responses; the shared contract
suite also exercises the new endpoint.

All 59 frontend tests passed, along with strict TypeScript checks and the
production build on Node 24.19.0. New tests verify that multiplication sends
numeric operands to `/api/multiply`, preserves the API result, displays the
result for the selected operation, and clears it on operation changes.

Browser checks against temporary backend/frontend servers verified `8 * 2 = 16`,
`-8 * 2 = -16`, `-8 * 0 = 0`, and presentation of `0.1 * 0.2` as `0.02`.
Overflow displayed the backend error; correcting the operand allowed successful
manual retry. The temporary servers were stopped after verification, and the
user's existing servers were not modified. Documentation was updated; division
remains pending user authorization.

## 14. Distinguish invalid format from an out-of-range number

After pasting a very long number in the multiplication form, I received the
generic invalid-number message. Explain this behavior, then distinguish invalid
numeric format from a number outside the supported range in the interface.

### AI implementation scope

Keep the approved finite numeric range and backend contract. Change the frontend
parser to return a valid numeric value or an explicit invalid-format/out-of-range
reason. Display a clear English message identifying the first invalid operand
field, with valid-format examples or a request for a smaller absolute value.
Continue blocking invalid input before fetch and clearing feedback on input
changes. Keep ordinary underflow to zero and calculation overflow behavior.
Add focused parsing and UI assertions; do not begin division in this refinement.

### AI review and validation outcome

All 68 frontend tests passed, along with strict TypeScript checks and production
build on Node 24.19.0. Tests distinguish malformed input from positive and negative
numeric overflow, a 400-digit pasted sequence, the maximum finite boundary, and
underflow to zero. UI cases verify both error reasons in both operand fields,
no request for invalid input, and clearing feedback after correction. Existing
calculation, API, loading, timeout, formatting, and retry tests also passed.

Browser checks confirmed the format message for `1,5`, the first-field range
message for `1e309`, and the second-field range message for `-1e309`.
No backend code or API contract changed. README and AGENTS.md were updated;
division remains pending authorization.

## 15. Division authorization

I tested the improved validation messages locally. Proceed with the next
requirement: division. Follow the agreed incremental implementation and review
process, then stop before beginning the next requirement.

### AI implementation scope

Add `POST /api/divide` computing `left / right`, preserving operand order.
Reject positive and negative zero divisors, including `0 / 0`, using HTTP 400
and the agreed `division_by_zero` code and `Cannot divide by zero.` message.
Keep finite-operand/result checks and the existing HTTP processing. Add Division
to the UI selector and shared API client. Display backend errors and allow manual
retry. Add relevant arithmetic, HTTP, and UI tests without new dependencies.

### AI review and validation outcome

Backend tests, build, static analysis, and formatting passed. Division tests
cover operand order, signs, a zero numerator, both zero divisor signs, `0 / 0`,
decimals, finite boundaries, small nonzero divisors, underflow, invalid operands,
and positive/negative overflow. Shared HTTP validation tests now run against all
four endpoints. Repeated success-response assertions were consolidated into a
test helper while preserving status, content type, exact JSON shape, and value
checks. Backend tests and checks passed again after this test-only cleanup.

All 72 frontend tests passed, along with strict TypeScript checks and the
production build on Node 24.19.0. New tests verify the division endpoint and
operand order, selected operation, clearing output, backend errors for `0` and
`-0` inputs, restored controls, and manual retry.

Browser checks against temporary backend/frontend servers verified `8 / 2 = 4`,
`2 / 8 = 0.25`, and presentation of `1 / 3` as `0.333333333333`. Positive
and negative zero inputs displayed `Cannot divide by zero.`, overflow displayed
the non-finite-result message, and corrected inputs allowed successful retry.
A direct HTTP request with raw JSON `-0` also returned the approved 400 error.
The temporary servers were stopped; the user's existing servers were not modified.
README and AGENTS.md were updated. All four arithmetic operations are implemented;
the remaining interface review, coverage reports, Docker, and final project review
remain pending their respective authorizations.

## 16. Interface, error handling, and responsiveness review authorization

The arithmetic implementation works after my review. Proceed with the next
requirement: review and complete any outstanding interface, error handling, and
basic mobile responsiveness requirements. Follow the incremental process and
stop for review before moving to coverage reports.

### AI review findings and implementation scope

The assistant reviewed the form, validation feedback, API and request states,
keyboard interaction, and responsive CSS. Manual checks found that client
validation left keyboard focus on Calculate and did not explicitly associate the
message with the invalid control. Disabling controls during a request also lost
keyboard focus after completion.

Client errors now identify their operand in structured UI feedback, focus it,
mark it with aria-invalid and a red border, and attach the message to its
accessible description. Input or operation changes clear the invalid mark and
message. General API errors do not mark valid fields invalid. A small focus
effect returns focus to Calculate after requests finish, without automatically
focusing it on initial render. Existing API behavior, numeric rules, layout, and
dependencies were retained.

### AI validation outcome

All 73 frontend tests passed, along with strict TypeScript checks and the
production build on Node 24.19.0. Existing validation cases were strengthened
to check focus, the invalid mark, accessible message association, and correction.
A new operation-change case verifies clearing field feedback. Success, API
failure, and timeout cases verify returned button focus. Existing request-state,
formatting, network, malformed-response, and retry assertions still passed.

Manual browser checks against the running backend verified Tab order, Enter
submission, focus on invalid fields, division-by-zero feedback, and retry.
Checks at 320x740, 375x812, 768x360, and 1280x800 showed the expected stacked or
side-by-side fields, readable long messages and scientific-notation results, and
no horizontal overflow. The viewport override was reset afterward. Physical
mobile keyboards and screen reader speech were not tested. No backend code
changed. Documentation was updated; coverage reports remain pending authorization.

## 17. Test review and coverage report authorization

Proceed with requirement 7: review whether the tests provide meaningful evidence
for the implemented requirements and generate frontend and backend coverage
reports. Follow the incremental process, document the results and remaining
coverage limits, and stop for review before beginning Docker.

### AI review findings and scope

The assistant reviewed the mathematical, HTTP, numeric input, API client, and
form tests against the approved behavior. A combined invalid error-envelope
case could short-circuit on its invalid code and fail to independently test
message validation. It was replaced with separate incorrect-type and blank
code/message cases. A JSON 503 response without the expected error envelope
was also checked. Backend tests now exercise a request-body transport read
failure through the real handler, expecting the agreed JSON invalid-input error
without disclosing the underlying transport error. No production arithmetic or
UI behavior changed.

The frontend adds @vitest/coverage-v8 5.0.3, matching Vitest, and an npm run
test:coverage command. V8 reports include application TypeScript files, including
unimported startup code, and exclude test files and test setup. Go reports include
all backend packages. No minimum percentage or new production failure mechanism
was introduced.

### AI validation outcome

All 77 frontend tests, strict TypeScript checks, and the production build passed
on Node 24.19.0. Backend tests with coverage, static analysis, and formatting
passed on Go 1.26.0. An isolated clean npm ci installation verified the updated
lockfile. HTML reports were generated and inspected in the browser; generated
coverage data and reports remain excluded from version control.

Backend statement coverage is 90.6% overall, 100% for the mathematical package,
and 96.2% for HTTP handling. Frontend statement coverage is 95.34% and branch
coverage is 94.82%. API access and numeric helpers have 100% on both; the form
has 100% statements and 95.83% branches. main.go and main.tsx have no unit-test
coverage and remain visible in the totals. Startup was exercised through prior
manual integration checks. Defensive exception/error paths remain uncovered and
are explained in the README instead of changing production code to raise coverage.

A focused audit introduced six selected regressions into temporary copies:
incorrect subtraction and multiplication operators, reversed division operands,
removal of zero-divisor rejection, and acceptance of blank API error codes or
messages. All six were detected by assertion failures, not compilation errors.
This is evidence for those specific assertions, not an exhaustive mutation score.
No mutation-testing dependency was added. Temporary copies and the local server
used to inspect reports were cleaned up. README and AGENTS.md were updated;
Docker remains pending separate authorization.

## 18. Docker authorization and backend image checkpoint

The coverage results look good. Proceed with Docker, following the agreed
incremental implementation and validation process.

### AI implementation scope

The assistant selected the backend Docker image as the first Docker checkpoint,
with user review before frontend/Nginx and Compose integration. Add a multi-stage
backend Dockerfile and a small build-context ignore file. Use a compatible Go
1.26 patch, run the existing backend tests in the build stage, compile with
CGO_ENABLED=0, and run the executable directly in a scratch image as a non-root
numeric user. Preserve the API contract and configurable PORT. Use host port
18080 in standalone instructions so the existing development API can keep using
8080. Do not change application behavior or scaffold frontend/Compose yet.

### AI validation outcome

The official golang:1.26.8-alpine3.24 image was downloaded and used successfully.
A build with --pull --no-cache passed, including all backend tests inside the
build container. Local static analysis and formatting also passed. Docker image
inspection confirmed a single runtime filesystem layer, a direct executable
entrypoint, and UID/GID 65532:65532. A process-inspection command was corrected
to include the PID column required by Docker top; the running UID was then
verified as 65532.

Fifteen HTTP checks against temporary containers passed: addition, subtraction
order, negative multiplication, division, positive/negative zero divisors,
overflow, invalid input, unknown path, incorrect method with Allow: POST,
missing/incompatible content type, the exact 1024-byte boundary, a 1025-byte
body, and a PORT=8081 override. Assertions checked statuses, JSON response
shapes, content types, result values, and error codes/messages.

The validated platform is Linux ARM64 using Docker Desktop 4.44.3 / Engine
28.3.2. An AMD64 build has not been verified. Temporary containers were stopped
and removed; the reusable backend image remains local. Existing development
servers were not changed. README and AGENTS.md were updated with verified
commands, image choices, port mappings, and limits. Frontend/Nginx and Compose
remain pending their subsequent checkpoints.

## 19. Frontend Docker and unified startup authorization

The Dockerfile currently covers only the backend. Add the frontend Dockerfile
now and provide one command that starts both frontend and backend together.

### AI implementation scope

Add a multi-stage frontend build using the approved Node 24 runtime and npm ci,
run the existing tests and production build, and serve compiled assets through
Nginx. Configure Nginx to proxy relative /api requests to the Go service while
preserving paths, status codes, JSON errors, and Allow. Add a small Docker Compose
configuration that builds and starts both services with docker compose up --build.
Publish the frontend at localhost:8081 by default, allow an APP_PORT override,
and keep backend port 8080 inside the Compose network so existing development
servers can coexist. Keep application code and arithmetic unchanged.

### AI review and validation outcome

Both images built successfully with --pull --no-cache. The frontend uses official
node:24.19.0-alpine3.24 and nginx:1.30.5-alpine3.24 images; npm ci, all 77 frontend
tests, strict TypeScript checks, and the production build passed in the container.
All backend tests also passed inside its build image. Compose configuration and
Nginx syntax checks passed. The frontend runtime contains compiled assets and
Nginx rather than the build-stage Node toolchain and dependencies.

Nginx streams API requests to Go and leaves body-size validation authoritative
there, preserving JSON 413 responses rather than replacing them with an HTML
proxy error. Paths are preserved through proxy_pass without a replacement URI.
Docker DNS and a resolving upstream follow backend service address changes.
Proxy retries are disabled. Compose orders service startup without claiming a
readiness guarantee; no health endpoint or wait script was added.

Compiled HTML, JavaScript, and CSS loaded. Seventeen API checks through Nginx
passed: four operations, decimal output, positive and negative zero divisors,
overflow, invalid input, unknown endpoint and API-prefix root, method and Allow,
missing/incompatible content type, exact/excessive body limits, and a body above
1 MiB. Inspection confirmed no backend host port was published. Browser checks
verified all operations, decimal presentation, division-by-zero feedback, manual
retry after input correction, an unavailable-service message with the backend
stopped, and successful manual retry after recreating that backend while keeping
the frontend running.

Validation used Linux ARM64 on Docker Desktop; AMD64 has not been verified.
README and AGENTS.md were updated with the unified startup command, browser
address, shutdown instructions, image versions, communication behavior, and
verification limits. An APP_PORT=8082 startup check passed, and docker compose down removed both
validation containers and their network. Reusable images remain local, and
existing native development servers were left running. The mandatory final
review remains the next checkpoint.

## 20. Approximate-result feedback approval

I tested very large operands and saw trailing digits replaced with zeros. Explain
whether this is mathematically exact and whether other calculators necessarily
behave the same way. I approve adding an approximately-equal symbol and a brief
English explanation so limited-precision results are transparent to the user.

### AI implementation scope and limits

The assistant explained that Number/float64 and the 12-significant-digit display
produce approximations, while calculators have different precision limits. Add
an approximation notice when presentation changes the backend value or either
operand is an integer outside JavaScript's safe integer range. The latter is a
conservative warning about potential input precision loss, even when that
particular large integer might be representable exactly. This does not prove
exactness when the marker is absent or detect every binary rounding/underflow
case. Keep the finite operand range, backend arithmetic, API values, and display
precision unchanged. Ordinary results such as 8 + 2 should remain without the
marker. Clear the note with stale output and associate it with the output's
accessible description.

### AI validation outcome

All 88 frontend tests and the coverage run passed on Node 24.19.0, along with
strict TypeScript checks and the production build. Added helper and UI cases
cover presentation rounding, repeating fractions, safe/unsafe integer boundaries,
both operand positions and signs, accessible feedback, the unchanged API payload,
and clearing on edits, operation changes, new requests, and request failure.
Frontend coverage is now 95.6% statements and 95.58% branches, retaining the
uncovered main.tsx startup in the totals. API and numeric helpers have 100%
statements and branches; App has 100% statements and 96.66% branches.

The updated frontend container also passed all 88 tests and its build. Browser
checks in an isolated Compose application on port 8082 verified 8 + 2 = 10
without a marker, approximate decimal feedback, and the user's large-integer
example with the approximately-equal symbol and explanatory note. A 320x740
viewport check found no horizontal overflow; the viewport was reset. Temporary
validation containers were removed, while the user's application on port 8081
was left running and requires a rebuild/refresh to load the change. README and
AGENTS.md were updated. The mandatory final review remains pending authorization.

## 21. UI subtitle removal

Remove the small Everyday math title from the UI.

The assistant removed the subtitle and its unused eyebrow CSS rule. All 24
form tests, strict TypeScript checks, and the production build passed. The
Docker frontend must be rebuilt to load this source change. Final review remains
pending.
