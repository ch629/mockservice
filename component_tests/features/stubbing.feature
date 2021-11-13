Feature: Stub Matching
    Scenario: A basic path equal to stub
        Given a definition is registered with payload
        """
        {
          "request": {
            "path": {
              "equal_to": "/abc"
            }
          },
          "response": {
            "body": {
              "foo": "bar"
            },
            "status": 200
          }
        }
        """
        When a request is sent with method "POST", path "/abc" and payload
        """
        {}
        """
        Then the response body should match
        """
        {
              "foo": "bar"
        }
        """
        And the response should have status code 200

    Scenario: No stub matched
        When a request is sent with method "POST", path "/abc" and payload
        """
        {}
        """
        Then the response should have status code 204
