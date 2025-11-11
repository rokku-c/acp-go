
<a href="https://agentclientprotocol.com/">
  <img alt="Agent Client Protocol" src="https://zed.dev/img/acp/banner-dark.webp">
</a>

# Agent Client Protocol (Go SDK)

This directory hosts the experimental Go SDK for the Agent Client Protocol (ACP). It is **not an official ACP implementation**; rather, it mirrors the protocol concepts demonstrated in the upstream ecosystem (see the Rust SDK for the reference design) so that we can integrate ACP inside our services.

## Project Goals

- Provide a production-ready Go binding that can talk to ACP-compliant clients and agents.
- Keep the protocol surface aligned with the public ACP schema and the official Rust SDK while adapting ergonomics to Go.

## Current Status

- Core request/response types, error helpers, and JSON-RPC framing are implemented in pure Go under `core/acp`.
- `AgentSideConnection` and `Client` abstractions support file system operations, terminal lifecycle management, and session notifications over ACP streams.
- Internal tests (`connection_test.go`, `client_inbound_test.go`) cover handshake flows, stream handling, and error propagation, but the surface is still evolving.

## Roadmap

- Expand coverage to the full ACP schema, including streaming updates and advanced terminal capabilities that remain stubbed today.
- Publish runnable examples that showcase agent-side and client-side integrations, following the layout used by the official Rust SDK.
- Harden interoperability by testing against upstream ACP reference clients/agents and capturing regressions in CI.
- Document contribution guidelines and release cadence once the API stabilises.

## Related Projects

- [Agent Client Protocol (Official Rust SDK)](https://github.com/agentclientprotocol/rust-sdk)
- [Agent Client Protocol Website](https://agentclientprotocol.com)


