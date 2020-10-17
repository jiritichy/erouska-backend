region                   = "europe-west1"
appengine_location       = "europe-west"
storage_location         = "EU"
cloudrun_location        = "europe-west1"
cloudscheduler_location  = "europe-west1"
kms_location             = "europe-west1"
network_location         = "europe-west1"
db_location              = "europe-west1"
db_name                  = "en-nform"
cloudsql_tier            = "db-custom-1-3840"
cloudsql_disk_size_gb    = "16"
generate_cron_schedule   = "*/15 * * * *"
cloudsql_max_connections = 10000
cloudsql_backup_location = "eu"
