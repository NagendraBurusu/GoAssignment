# Go Assignment

## Nagnedra Burusu
- nagendra.burusu@thoughtclan.com

### Step by step of Building Students App
- Loading the  configuration from the .env file.And loading before  main method by using init() method.
- Making Mysql connection by using these configuration.
- Adding the crud operations and alive and Ready check functions to the student app.
- Added middilewares and provided Authentication API to generate JWT token by taking userId and password as a String
- Used middleware to check the authentication as well as provide the user ionformation via contex to db layer
- updated the create and update method to take the user information and save the changes. 