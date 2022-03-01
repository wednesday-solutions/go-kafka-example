# Producer


The producer publishes to 2 topics
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

