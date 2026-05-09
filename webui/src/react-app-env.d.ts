/// <reference types="react-scripts" />

declare module "*.css" {
  const content: Record<string, string>;
  export default content;
}

declare module "babel-plugin-relay/macro" {
  export { graphql as default } from "react-relay";
}