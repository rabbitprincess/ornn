table "newtable" {
  schema = schema.public
  column "a" {
    null = true
    type = character_varying
  }
  column "b" {
    null = true
    type = character_varying
  }
  column "seq" {
    null = false
    type = bigint
  }
  primary_key {
    columns = [column.seq]
  }
}
schema "public" {
}
