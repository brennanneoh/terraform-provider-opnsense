# Import the existing wan assignment before managing it:
#   terraform import opnsense_interfaces_assignment.wan wan
resource "opnsense_interfaces_assignment" "wan" {
  device      = "vtnet0"
  description = "WAN"
}

# Assigning a device to a new logical name creates an optional interface
# (opt1, opt2, ...); it does not create "wan" or "lan".
resource "opnsense_interfaces_assignment" "opt_example" {
  device      = "vtnet2"
  description = "opt example"
}
