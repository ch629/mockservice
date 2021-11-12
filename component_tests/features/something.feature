Feature: Testing
    Scenario: Foo
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
