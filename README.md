# Gamma, by [Teleport](https://goteleport.com)

_**G**ithub **A**ctions **M**onorepo **M**agic **A**utomation_

Gamma is a tool that sets out to solve a few shortcomings when it comes to managing and maintaining multiple GitHub actions.

## What does it do?

- 🚀 No more including the compiled source code in your commits
- 🚀 Automatically build all your actions into individual, publishable repos
- 🚀 Share schema definitions between actions
- 🚀 Version all actions separately

Gamma allows you to have a monorepo of actions that are then built and deployed into individual repos. Having each action in its own repo allows for the action to be published on the Github Marketplace.

Gamma also goes further when it comes to sharing common `action.yml` attributes between actions. Actions in your monorepo can extend upon other YAML files and bring in their `inputs`, `branding`, etc - reducing code duplication and making things easier to maintain.

## How to use

This assumes you're using `yarn` with workspaces. Each workspace is an action.

Your root `package.json` should look like:

```json
{
  "name": "actions-monorepo",
  "private": true,
  "workspaces": [
    "actions/*"
  ]
}
```

Each action then lives under the `actions/` directory.

Each action should be able to be built via `yarn build`. We recommend [ncc](https://github.com/vercel/ncc) for building your actions. The compiled source code should end up in a `dist` folder, relative to the action. You should add `dist/` to your `.gitignore`.

`actions/example/package.json`

```json
{
  "name": "example",
  "version": "1.0.0",
  "repository": "https://github.com/mono-actions/example.git",
  "scripts": {
    "build": "ncc build ./src/index.ts -o dist"
  },
  "dependencies": {
    "@actions/core": "^1.10.0"
  },
  "devDependencies": {
    "@types/node": "^18.8.2",
    "@vercel/ncc": "^0.34.0",
    "typescript": "^4.8.4"
  }
}
```

The `repository` field is where the compiled action will deployed to.

`actions/example/action.yml`

This is where Gamma can really shine. You can define your `action.yml` as normal, whilst also extending on other YAML files for common attributes.

```yaml
name: Example Action
description: This is an example action
extend:
  - from: '@/shared/common.yml'
    include:
      - field: inputs
        include:
          - version
      - field: runs
      - field: author
      - field: branding
```

`@/` refers to the root of the directory. `@/shared/common.yml` would resolve to `shared/common.yml`, which can look like this:

`shared/common.yml`

```yaml
author: Gravitational, Inc.
inputs:
  version:
    required: true
    description: 'Specify the version without the preceding "v"'
branding:
  icon: terminal
  color: purple
runs:
  using: node16
  main: dist/index.js
```

Gamma will compile this and publish the final `action.yml` to the correct repository.

`github.com/mono-actions/example/action.yml`

```yaml
name: Example Action
description: This is an example action
author: Gravitational Inc.
inputs:
    version:
        description: Specify the version without the preceding "v"
        required: true
runs:
    using: node16
    main: dist/index.js
branding:
    icon: terminal
    color: purple
```

The built source code will also be committed, so you end up with a publishable Github Action.

## Use in GitHub actions

You can use this in your GitHub action workflows via [setup-gamma](https://github.com/gravitational/setup-gamma).
