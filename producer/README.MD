# Producer

This application publishes messages to 2 topic
- [issued-token](https://github.com/wednesday-solutions/go-kafka-example/blob/main/producer/resolver/auth_mutations.resolvers.go#L66)
- [new-user-created](https://github.com/wednesday-solutions/go-kafka-example/blob/main/producer/resolver/user_mutations.resolvers.go#L57)

- issued-token
  To publish a message on this topic simply login
  ```
    mutation login {
      login(username:"admin", password:"adminuser") {
        token
        refreshToken
      }
    }
  ```
- new-user-created
  To publish a message on this topic create a user
    ```
    mutation CreateUser {
      createUser(input:{
        firstName:"Mohammed Ali"
        lastName:"Chherawalla",
        username:"mohammedali",
        password:"Kasd123EHU",
        email:"mac@wednesday.is",
        roleId:5
      }) {
        user{
          id
          firstName
          lastName
          username
          email
        }
      }
    }
    ```

