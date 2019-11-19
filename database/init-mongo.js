db.createUser(
  {
    user: "izzaturrahman19",
    pwd: "mindtrex",
    roles: [
       { role: "readWrite", db: "mindtrex" }
    ]
  }
)