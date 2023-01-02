table "test" {
  schema  = schema.test
  charset = "utf8mb4"
  collate = "utf8mb4_general_ci"
  column "id" {
    null           = false
    type           = int
    unsigned       = true
    auto_increment = true
  }
  primary_key {
    columns = [column.id]
  }
}
schema "test" {
  charset = "latin1"
  collate = "latin1_swedish_ci"
}
