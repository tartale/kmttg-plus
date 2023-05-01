import { Episode, Movie, Series, Show, ShowKind } from "./ShowListing";


export const startJob = () => {
  fetch("http://localhost:8181/getMyShows?limit=50&tivo=Living%20Room&offset=0", {
    "credentials": "omit",
    "headers": {
      "Accept": "application/json, text/javascript, */*; q=0.01",
      "Accept-Language": "en-US,en;q=0.5",
      "X-Requested-With": "XMLHttpRequest",
      "Sec-Fetch-Dest": "empty",
      "Sec-Fetch-Mode": "cors",
      "Sec-Fetch-Site": "same-origin",
      "Sec-GPC": "1",
      "Pragma": "no-cache",
      "Cache-Control": "no-cache",
    },
    "method": "GET",
    "mode": "cors"
  })
    .then((response) => response.json())
    .then((jsonArray) => {
    })
    .catch((error) => console.error(error));

}

export const waitForJob = () => {
  fetch("http://localhost:8181/getMyShows?limit=50&tivo=Living%20Room&offset=0", {
    "credentials": "omit",
    "headers": {
      "Accept": "application/json, text/javascript, */*; q=0.01",
      "Accept-Language": "en-US,en;q=0.5",
      "X-Requested-With": "XMLHttpRequest",
      "Sec-Fetch-Dest": "empty",
      "Sec-Fetch-Mode": "cors",
      "Sec-Fetch-Site": "same-origin",
      "Sec-GPC": "1",
      "Pragma": "no-cache",
      "Cache-Control": "no-cache",
    },
    "method": "GET",
    "mode": "cors"
  })
    .then((response) => response.json())
    .then((jsonArray) => {
    })
    .catch((error) => console.error(error));

}

export const downloadShow = () => {
  fetch("http://localhost:8181/getMyShows?limit=50&tivo=Living%20Room&offset=0", {
    "credentials": "omit",
    "headers": {
      "Accept": "application/json, text/javascript, */*; q=0.01",
      "Accept-Language": "en-US,en;q=0.5",
      "X-Requested-With": "XMLHttpRequest",
      "Sec-Fetch-Dest": "empty",
      "Sec-Fetch-Mode": "cors",
      "Sec-Fetch-Site": "same-origin",
      "Sec-GPC": "1",
      "Pragma": "no-cache",
      "Cache-Control": "no-cache",
    },
    "method": "GET",
    "mode": "cors"
  })
    .then((response) => response.json())
    .then((jsonArray) => {
    })
    .catch((error) => console.error(error));

}
