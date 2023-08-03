import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import { useEffect, useState } from "react";
import "./ShowListing.css";
import { ShowHeader, ShowRow } from "./ShowRow";
import { Show } from "../services/generated/graphql-types"
import "./TivoStyle.css";
import {
  useQuery,
  gql
 } from "@apollo/client";

export type ShowSortField = 'kind' | 'title' | 'recordedOn';

const GET_RECORDINGS = gql`
 query getRecordings {
  tivos {
    recordings(limit:50) {
      kind
      recordingID
      title
      recordedOn
      ... on Series {
        episodes {
          kind
          recordingID
          episodeTitle
          seasonNumber
          episodeNumber
          episodeDescription
        }
      }
    }
  }
}`;

function ShowListingComponent(props: any) {
  const { showlisting } = props
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(() => setShows(showlisting), [showlisting]);

  return (
    <TableContainer
      component={Paper}
      sx={{ background: "linear-gradient(to bottom, #162c4f, #000000);" }}
      {...props}
    >
      <Table className="showListingTable">
        <ShowHeader/>
        <TableBody>
          {shows.map((show) => (
            <ShowRow key={show.recordingID} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default function ShowListing(props: any) {
  const { loading, error, data } = useQuery(GET_RECORDINGS);

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    console.error(error);
    return <div>Error!</div>;
  }

  const showListing: Show[] = data.tivos[0].recordings
  return <ShowListingComponent showlisting={showListing} {...props}/>;
};
