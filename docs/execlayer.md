# Execution layer configuration
The execution layer configuration contains two pieces.  The first, which is mandatory, provides the fee recipient to the execution layer.  The second, which is optional, provides access to MEV (maximum extractable value) relays.

## Basic configuration

Configuration is provided under the `blockrelay` key.  The minimal configuration is as follows:

```YAML
blockrelay:
  fallback-fee-recipient: '0x0123…cdef'
```

This will use the provided address `0x0123…cdef` as the fee recipient for all beacon block proposals.  Block proposals requested from beacon nodes will use their local execution client to obtain execution payloads.

For more advanced configurations an execution configuration file is required.  Access to the configuration file is usually through a simple URL, for example:

```YAML
blockrelay:
  fallback-fee-recipient: '0x0123…cdef'
  config:
    url: 'file:///home/vouch/config.json'
```

(Note that the `fallback-fee-recipient` value is still present, and is required for Vouch to operate in the situation that the configuration file is inaccessible or unreadable.)

The URL can be any valid [Majordomo](https://github.com/wealdtech/go-majordomo) URL.  Of special note is the [HTTP confidant](https://github.com/wealdtech/go-majordomo/blob/master/confidants/http/service.go#L32), for which Vouch provides additional features.  Notably, the call is made as a POST request with its body containing the public keys of the validators for which Vouch is currently validating, for example:

```json
[
  "0x1111…1111",
  "0x2222…2222"
]
```

Also, additional configuration parameters can be provided to secure the connection using HTTPS.  A full example of such a configuration is:

```YAML
blockrelay:
  fallback-fee-recipient: '0x0123…cdef'
  config:
    url: 'https://www.example.com/config.json'
    client-cert: 'file:///home/vouch/certs/my-client.crt'
    client-key: 'file:///home/vouch/certs/my-client.key'
    ca-cert: 'file:///home/vouch/certs/ca.crt'
```

The combination of the list of public keys and certificate-level authentication allows for servers to provide dynamic execution configuration information.

The execution configuration file referenced by the configuration URL allows for a default configuration alongside per-validator overrides, for example:

```json
{
  "proposer_config": {
    "0xaaaa…aaaa": {
      "fee_recipient": "0x1111…1111"
    },
    "0xbbbb…bbbb": {
      "fee_recipient": "0x2222…2222"
    }
  },
  "default_config": {
    "fee_recipient": "0x0123…cdef"
  }
}
```

In this configuration the validator with the public key `0xaaaa…aaaa` will use the address `0x1111…1111` as its fee recipient, the validator with the public key `0xbbbb…bbbb` will use the address `0x2222…2222` as its fee recipient, and all other validators will use the address `0x0123…cdef` as their fee recipients.

## MEV configuration

Vouch acts as an MEV-boost server, talking to MEV relays and accepting requests from beacon nodes for execution payload headers.

```json
{
  "proposer_config": {
    "0xaaaa…aaaa": {
      "fee_recipient": "0x1111…1111",
      "builder": {
        "enabled": true,
        "relays": [
          "relay1.example.com",
          "relay2.example.com"
        ]
      }
    },
    "0xbbbb…bbbb": {
      "fee_recipient": "0x2222…2222",
      "builder": {
        "enabled": false
      }
    }
  },
  "default_config": {
  "fee_recipient": "0x0123…cdef",
    "builder": {
      "enabled": true,
      "relays": [
        "relay1.example.com",
        "relay2.example.com"
      ]
    }
  }
}
```

## Beacon node configuration
By default, Vouch's MEV-boost service listens on port 18550.  Beacon nodes used by Vouch should be configured to talk directly to this, rather than other MEV-boost services or directly to relays.  Details of how to configure the beacon nodes is listed in the relevant client's documentation.  If a different port is required this can be set in the block relay configuration, for example:

```YAML
blockrelay:
  listen-address: '0.0.0.0:12345'
  ...
```

## Gas limit
Block proposers have the ability to alter the gas limit as part of their block proposal process.  In general it is recommended that this value be left to the Vouch default, as it requires a majority of block proposers to agree on a new value for it to be reached, however if there is a requirement to change this then it can be done in the execution configuration file.  A sample execution configuration file that includes changing gas limit to 100000000 is shown below:

```json
{
  "proposer_config": {
    "0xaaaa…aaaa": {
      "fee_recipient": "0x1111…1111",
      "gas_limit": 100000000,
      "builder": {
        "enabled": true,
        "relays": [
          "relay1.example.com",
          "relay2.example.com"
        ]
      }
    },
    "0xbbbb…bbbb": {
      "fee_recipient": "0x2222…2222",
      "gas_limit": 100000000,
      "builder": {
        "enabled": false
      }
    }
  },
  "default_config": {
  "fee_recipient": "0x0123…cdef",
  "gas_limit": 100000000,
    "builder": {
      "enabled": true,
      "relays": [
        "relay1.example.com",
        "relay2.example.com"
      ]
    }
  }
}
```

## Dynamic configuration file

The execution configuration file is re-read each epoch, which allows for changes to take place without restarting Vouch.
