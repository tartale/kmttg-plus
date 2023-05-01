/**
 * @generated SignedSource<<602239fe8ab5abdcb02be1dc45e6ab17>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest, Query } from 'relay-runtime';
export type TivoSelectorQuery$variables = {};
export type TivoSelectorQuery$data = {
  readonly tivos: ReadonlyArray<{
    readonly name: string;
  }>;
};
export type TivoSelectorQuery = {
  response: TivoSelectorQuery$data;
  variables: TivoSelectorQuery$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "concreteType": "Tivo",
    "kind": "LinkedField",
    "name": "tivos",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "name",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "TivoSelectorQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "TivoSelectorQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "3eb7ee3bacfa853287ee1a033ebb2f7e",
    "id": null,
    "metadata": {},
    "name": "TivoSelectorQuery",
    "operationKind": "query",
    "text": "query TivoSelectorQuery {\n  tivos {\n    name\n  }\n}\n"
  }
};
})();

(node as any).hash = "71644f48fb85f1420f55ac3fdd5177f8";

export default node;
