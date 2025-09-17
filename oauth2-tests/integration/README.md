# Integration Tests for OAuth2 Server

This directory contains integration tests for the Go OAuth2 server, written using `pytest-bdd`. These tests simulate client registration flows, including successful and various failure scenarios, to ensure the server's API behaves as expected.

## Prerequisites

Before running the tests, ensure you have:

*   **Python 3.8+** installed.
*   **pip** (Python package installer).
*   The **Go OAuth2 server** running on `http://localhost:8080` (or adjust the `BASE_URL` in `tests/test_register_client_steps.py` accordingly).

## Installation

1.  Navigate to this directory in your terminal:
    ```bash
    cd oauth2-tests/integration
    ```
2.  **Create and Activate a Virtual Environment (Recommended)**:
    ```bash
    python3 -m venv venv
    source venv/bin/activate
    ```
    *   On Windows, use `.\venv\Scripts\activate`.

3.  Install the required Python dependencies:
    ```bash
    pip install -r requirements.txt
    ```

## Running Tests

To execute the integration tests, use `pytest`.

### Basic Run

To run all tests with verbose output:
```bash
pytest -v -s tests/test_register_client_steps.py
```
*   `-v`: Enables verbose output.
*   `-s`: Disables stdout/stderr capture, allowing `print()` statements in your test steps to be visible.

### Generating an HTML Report

To generate a detailed HTML report of the test results, use the `pytest-html` plugin:
```bash
pytest --html=report.html --self-contained-html tests/test_register_client_steps.py
```
This will create a `report.html` file in the current directory.

## Viewing the Test Report

After generating the HTML report, you can open it in your web browser:
```bash
open report.html
# On Linux, you might use: xdg-open report.html
# On Windows, you might use: start report.html
```

## Test Structure

*   `features/`: Contains Gherkin feature files (`.feature`) that describe the test scenarios in a human-readable format.
*   `tests/`: Contains Python step definitions (`.py`) that implement the logic for the Gherkin steps.
*   `requirements.txt`: Lists all Python package dependencies.
*   `.gitignore`: Specifies files and directories to be ignored by Git.

---
