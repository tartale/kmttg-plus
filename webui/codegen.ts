import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  overwrite: true,
  schema: "../go/api/**/*.graphql",
  // documents: "src/services/**/*.graphql",
  watchConfig: {
    usePolling: true,
    interval: 1000,
  },
  generates: {
    './src/services/generated/graphql-types.ts': {
      plugins: ['typescript'],
      config: {
        scalars: {
          Date: {
            input: "string",
            output: "Date",       
          },
          Time: {
            input: "string",
            output: "Date"
          }
        },
      }
    }
  }
}
export default config
