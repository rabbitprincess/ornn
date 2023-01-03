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
}
schema "public" {
}
