Feature: Register a new OAuth2 client
  Scenario: Successfully register a new client
    Given I have a new client registration payload
    When I send the registration request
    Then I should receive a 201 status code
    And the response should contain a client_id and client_secret

  Scenario: Fail to register client with missing client_name
    Given I have a client registration payload missing "client_name"
    When I send the registration request
    Then I should receive a 400 status code
    And the response should contain an error "invalid_request"
    And the response should contain an error description "The request is missing a required parameter, includes an invalid parameter value, or is otherwise malformed."

  Scenario: Fail to register client with an invalid redirect_uri
    Given I have a client registration payload with invalid "redirect_uris"
    When I send the registration request
    Then I should receive a 400 status code
    And the response should contain an error "invalid_redirect_uri"
    And the response should contain an error description "One or more redirect URIs are invalid or missing"

  Scenario: Fail to register client with an invalid scope
    Given I have a client registration payload with invalid "scope"
    When I send the registration request
    Then I should receive a 400 status code
    And the response should contain an error "invalid_scope"
    And the response should contain an error description "The requested scope is invalid, unknown, or malformed."
