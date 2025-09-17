import pytest
import requests
import uuid
import json # Import json for pretty printing
from pytest_bdd import scenarios, given, when, then, parsers

BASE_URL = "http://localhost:8080"  # Replace with your OAuth2 server base URL

scenarios("../features/register_client.feature")


@pytest.fixture
@given("I have a new client registration payload")
def base_client_payload():
    """Initializes a base client registration payload with unique client_name."""
    payload = {
        "client_name": f"test-client-app-{uuid.uuid4()}",
        "grant_types": ["authorization_code", "client_credentials"],
        "response_types": ["code"],
        "token_endpoint_auth_method": "client_secret_basic",
        "redirect_uris": ["http://localhost:8080/callback"],
        "scope": "email",
    }
    return payload


@pytest.fixture
@given("client_payload")
def client_payload(base_client_payload):
    """Provides a mutable copy of the base client payload for each scenario."""
    return base_client_payload.copy()


@given(parsers.parse('I have a client registration payload missing "{field_name}"'))
def client_payload_missing_field(client_payload, field_name):
    """Modifies the client payload by removing a specified field."""
    if field_name in client_payload:
        del client_payload[field_name]
    # Ensure redirect_uris is still valid for other tests if client_name is missing
    if field_name == "client_name":
        client_payload["redirect_uris"] = ["http://localhost:8080/callback"]
        client_payload["scope"] = "openid"
    return client_payload


@given(
    parsers.parse('I have a client registration payload with invalid "redirect_uris"')
)
def client_payload_invalid_redirect_uri(client_payload):
    """Modifies the client payload to include an invalid redirect_uri."""
    client_payload["redirect_uris"] = ["invalid-uri"]
    return client_payload


@given(parsers.parse('I have a client registration payload with invalid "scope"'))
def client_payload_invalid_scope(client_payload):
    """Modifies the client payload to include an invalid scope."""
    client_payload["scope"] = "invalid_scope_name"
    return client_payload


@when("I send the registration request")
def send_registration_request(client_payload):
    """Sends the client registration request to the OAuth2 server."""
    print(f"\n[When] Sending registration request with payload:\n{json.dumps(client_payload, indent=2)}")
    headers = {"Content-Type": "application/json"}
    response = requests.post(
        f"{BASE_URL}/oauth/register", json=client_payload, headers=headers
    )
    client_payload["response"] = response
    print(f"[When] Received response (status {response.status_code}):\n{response.text}")
    return client_payload


@then("I should receive a 201 status code")
def check_status(client_payload):
    """Asserts that the HTTP response status code is 201 (Created)."""
    response = client_payload["response"]
    assert response.status_code == 201, (
        f"Expected 201, but got {response.status_code}. Response: {response.text}"
    )


@then("I should receive a 400 status code")
def check_bad_request_status(client_payload):
    """Asserts that the HTTP response status code is 400 (Bad Request)."""
    response = client_payload["response"]
    assert response.status_code == 400, (
        f"Expected 400, but got {response.status_code}. Response: {response.text}"
    )


@then("the response should contain a client_id and client_secret")
def check_client_credentials(client_payload):
    """Asserts that the response JSON contains 'client_id' and 'client_secret'."""
    try:
        json_data = client_payload["response"].json()
    except json.JSONDecodeError:
        pytest.fail(f"Response is not valid JSON: {client_payload['response'].text}")
    assert "client_id" in json_data, f"'client_id' not found in response: {json_data}"
    assert "client_secret" in json_data, f"'client_secret' not found in response: {json_data}"


@then(parsers.parse('the response should contain an error "{error_code}"'))
def check_error_code(client_payload, error_code):
    """Asserts that the response JSON contains the specified error code."""
    try:
        json_data = client_payload["response"].json()
    except json.JSONDecodeError:
        pytest.fail(f"Response is not valid JSON: {client_payload['response'].text}")
    assert "error" in json_data, f"'error' field not found in response: {json_data}"
    assert json_data["error"] == error_code, (
        f"Expected error code '{error_code}', but got '{json_data.get('error')}'. Full response: {json_data}"
    )


@then(
    parsers.parse(
        'the response should contain an error description "{error_description}"'
    )
)
def check_error_description(client_payload, error_description):
    """Asserts that the response JSON contains the specified error description."""
    try:
        json_data = client_payload["response"].json()
    except json.JSONDecodeError:
        pytest.fail(f"Response is not valid JSON: {client_payload['response'].text}")
    assert "error_description" in json_data, f"'error_description' field not found in response: {json_data}"
    assert json_data["error_description"] == error_description, (
        f"Expected error description '{error_description}', but got '{json_data.get('error_description')}'. Full response: {json_data}"
    )
