job "g2ww" {
  datacenters = ["dc1"]

  group "alert" {
    count = 1

    network {
      dns {
        servers = ["114.114.114.114", "8.8.8.8", "8.8.4.4"]
      }

      port "http" {
        static = 2408
        to     = 2408
      }
    }

    service {
      name = "g2ww"
      port = "http"

      check {
        type     = "tcp"
        port     = "http"
        interval = "30s"
        timeout  = "3s"
      }
    }

    task "g2ww" {
      env {
        PORT    = "${NOMAD_PORT_http}"
        NODE_IP = "${NOMAD_IP_http}"
      }

      driver = "podman"

      config {
        image = "80x86/g2ww:latest"
        ports = ["http"]
      }
    }
  }
}
