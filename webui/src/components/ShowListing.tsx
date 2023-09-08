import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import { useEffect, useState, useRef } from "react";
import "./ShowListing.css";
import { ShowHeader, ShowRow } from "./ShowRow";
import { Show } from "../services/generated/graphql-types"
import "./TivoStyle.css";
import { useQuery, gql } from "@apollo/client";

export type ShowSortField = 'kind' | 'title' | 'recordedOn';

const GET_RECORDINGS = gql`
 query getRecordings($offset: Int, $limit: Int) {
  tivos {
    name
    shows(offset: $offset, limit: $limit) {
      id
      kind
      title
      description
      recordedOn
      ... on Movie {
        imageURL(height: 512, width: 512)
      }
      ... on Series {
        episodes {
          id
          kind
          episodeTitle
          seasonNumber
          episodeNumber
          episodeDescription
        }
        imageURL(height: 512, width: 512)
      }
    }
  }
}`;

function ShowListingComponent(props: any) {
  const { showListing, loadMoreData, isLoadingMore, ...remainingProps } = props;
  const [shows, setShows] = useState<Show[]>([]);
  const lastRowRef = useRef<HTMLTableRowElement>(null);

  useEffect(() => setShows(showListing), [showListing]);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !isLoadingMore) {
          loadMoreData();
        }
      },
      {
        threshold: 0.1,
      }
    );

    if (lastRowRef.current) {
      observer.observe(lastRowRef.current);
    }

    return () => {
      if (lastRowRef.current) {
        observer.unobserve(lastRowRef.current);
      }
    };
  }, [isLoadingMore]);

  return (
    <TableContainer
      component={Paper}
      sx={{ background: 'linear-gradient(to bottom, #162c4f, #000000);' }}
      {...remainingProps}
    >
      <Table className="showListingTable">
        <ShowHeader />
        <TableBody>
          {shows.map((show, index) => (
            <ShowRow key={show.id} show={show} />
          ))}
          {isLoadingMore && (
            <TableRow ref={lastRowRef} key="lastRowRef">
              <TableCell colSpan={6}>
                {'Loading...'}
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

export default function ShowListing(props: any) {
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const { loading, error, data, fetchMore } = useQuery(GET_RECORDINGS, {
    variables: {
      offset: 0,
      limit: 50,
    },
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    console.error(error);
    return <div>Error!</div>;
  }

  const showListing: Show[] = data.tivos[0].shows;

  const loadMoreData = () => {
    setIsLoadingMore(true);

    fetchMore({
      variables: {
        offset: showListing.length,
      },
      updateQuery: (prev, { fetchMoreResult }) => {
        setIsLoadingMore(false);
        if (!fetchMoreResult) return prev;
        return {
          tivos: [...prev.tivos, ...fetchMoreResult.tivos],
        };
      },
    });
  };

  return (
    <div>
      <ShowListingComponent showListing={showListing} {...props} loadMoreData={loadMoreData} isLoadingMore={isLoadingMore} />
    </div>
  );
}
