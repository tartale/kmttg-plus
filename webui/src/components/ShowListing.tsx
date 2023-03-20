import React from "react";
import "./ShowListing.css";

export interface Show {
  id?: number;
  recordedOn: string;
  title: string;
  episodeTitle: string;
}

export function ShowListing(props: {shows: Show[]}) {
  const {shows} = props;
  return (
    <table className="table">
      {/* <thead>
        <tr>
          <th>Recorded On</th>
          <th>Title</th>
          <th>Episode Title</th>
        </tr>
      </thead> */}
      <tbody>
        {shows.map((show: Show, index: number) => (
          <tr key={index}>
            <td>{show.recordedOn}</td>
            <td>{show.title}</td>
            <td>{show.episodeTitle}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};
