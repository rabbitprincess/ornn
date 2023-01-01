table "user" {
  schema = schema.test
  column "seq" {
    null           = false
    type           = bigint
    auto_increment = true
  }
  column "id" {
    null    = false
    type    = varchar(50)
    default = ""
  }
  column "name" {
    null    = false
    type    = varchar(50)
    default = ""
  }
  primary_key {
    columns = [column.seq]
  }
}
schema "test" {
  charset = "utf8mb4"
  collate = "utf8mb4_general_ci"
}
