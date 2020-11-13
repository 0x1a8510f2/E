# E - Matrix Command & Control
<p align="center">
  <img alt="E" src="https://user-images.githubusercontent.com/37966924/98878631-588ac980-247b-11eb-897c-7c7e2c8ad1ad.png" />
</p>
<p align="center">
  <a href="https://github.com/TR-SLimey/E/actions?query=workflow%3ADockerBuild">
    <img alt="Docker build status" src="https://github.com/TR-SLimey/E/workflows/DockerBuild/badge.svg" />
  </a>
  <a href="https://matrix.to/#/#e:an-atom-in.space">
    <img alt="Matrix" src="https://img.shields.io/matrix/e:an-atom-in.space?server_fqdn=matrix-client.matrix.org" />
  </a>
  <a href="https://github.com/TR-SLimey/E/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues-raw/TR-SLimey/E" />
  </a>
  <a href="https://github.com/TR-SLimey/E/issues">
    <img alt="GitHub issues" src="https://img.shields.io/github/issues-closed-raw/TR-SLimey/E" />
  </a>
  <a href="https://github.com/TR-SLimey/E/blob/master/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/TR-SLimey/E" />
  </a>
  <a href="https://github.com/TR-SLimey/E/pulls">
    <img alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr-raw/TR-SLimey/E" />
  </a>
  <a href="https://github.com/TR-SLimey/E/commits">
    <img alt="GitHub commit activity" src="https://img.shields.io/github/commit-activity/m/TR-SLimey/E" />
  </a>
  <a href="https://hub.docker.com/repository/docker/trslimey/e_mc2">
    <img alt="Docker Stars" src="https://img.shields.io/docker/stars/trslimey/e_mc2" />
  </a>
</p>
<br/>

## What is E? 🤨
<table border="0px">
  <tr>
    <td valign="top">
      <pre>
        <code>
  E = mc²
     /  \
 matrix  cc
        /  \
  command & control
        </code>
      </pre>
    </td>
    <td valign="top">
      <p>
        E is a multi-purpose and multi-protocol command & control server designed to receive commands through the <a href="https://matrix.org">Matrix</a> chat protocol and forward them to clients. It is also specifically designed to be versatile in that it can be used to bridge any data between Matrix and other protocols as long as valid <code>esockets</code> are implemented for that protocol.
    </td>
  </tr>
</table>

## Why is E? 🤔
While C&C servers are usually thought of in the context of malware, and E could indeed be used for this purpose (in which case I accept no responsibility for what you do), E can also be used for any number of things, such as:
- controlling IoT devices via a text message - Matrix is already great at that, but with E, a custom protocol can be implemented to fit specific requirements
- administering a large number of devices in business settings - E could translate between the protocols of tools you already use and Matrix 
- sending Matrix messages over any protocol providing there is an esocket for it
- proxying Matrix traffic, if you want to do that for some reason instead of hosting a homeserver (requires a Matrix CS API esocket)

## When is E? 📅🤔
What??

## How Tos 🔍
```
// TODO
```
