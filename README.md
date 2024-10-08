# BDD demo
A simple demo for BDD testing in Golang.

Reference:
- https://medium.com/tiket-com/go-with-cucumber-an-introduction-for-bdd-style-integration-testing-7aca2f2879e4

Run the tests:
```
 go test -v -timeout 30s -run ^TestFeatures$ github.com/hedon-go-road/bdd-demo/tests
```
output:
```bash
=== RUN   TestFeatures
Feature: Book management
  In order to use book API
  As a Librarian
  I need to be able to manage books.
=== RUN   TestFeatures/then_user_try_to_insert_one_book,_one_created_book_should_be_displayed_by_the_system
2024/09/19 19:25:42 connected
2024/09/19 19:25:42 running migrations

  Scenario: then user try to insert one book, one created book should be displayed by the system # features/book.feature:6
    When I send "POST" request to "/books" with payload:                                         # <autogenerated>:1 -> *apiFeature
      """
      {
        "id": 1,
        "title": "Dune",
        "author": "Frank Herbert"
      }
      """
    Then the response code should be 201                                                         # <autogenerated>:1 -> *apiFeature
    And the response payload should match json:                                                  # <autogenerated>:1 -> *apiFeature
      """
      [
        {
          "id": 1,
          "title": "Dune",
          "author": "Frank Herbert"
        }
      ]
      """

1 scenarios (1 passed)
3 steps (3 passed)
44.807417ms
--- PASS: TestFeatures (0.05s)
    --- PASS: TestFeatures/then_user_try_to_insert_one_book,_one_created_book_should_be_displayed_by_the_system (0.04s)
PASS
ok      github.com/hedon-go-road/bdd-demo/tests 1.912s
```