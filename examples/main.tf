module "this" {
  source         = "../"
  environment    = "staging"
  service_name   = "spot-interruption-handler"
  description    = "Lambda to handle spot instance interruption"
}
