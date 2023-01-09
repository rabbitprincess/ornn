table "tbltest" {
  schema  = schema.test
  charset = "utf8mb4"
  collate = "utf8mb4_general_ci"
  column "seq" {
    null           = false
    type           = int
    unsigned       = true
    auto_increment = true
  }
  column "id2" {
    null    = true
    type    = varchar(50)
    default = ""
  }
  column "address" {
    null    = true
    type    = varbinary(50)
    default = ""
  }
  column "registered" {
    null = false
    type = bool
  }
  primary_key {
    columns = [column.seq]
  }
}
table "user" {
  schema  = schema.test
  charset = "utf8mb4"
  collate = "utf8mb4_general_ci"
  column "seq" {
    null           = false
    type           = int
    unsigned       = true
    auto_increment = true
  }
  column "id2" {
    null    = true
    type    = varchar(50)
    default = ""
  }
  column "address" {
    null    = true
    type    = varbinary(50)
    default = ""
  }
  column "registered" {
    null = false
    type = bool
  }
  primary_key {
    columns = [column.seq]
  }
}
schema "test" {
  charset = "latin1"
  collate = "latin1_swedish_ci"
}
