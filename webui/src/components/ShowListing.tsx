import React, { useState, useEffect } from "react";
import "./TivoStyle.css";
import "./ShowListing.css";

export interface Show {
  id?: number;
  title: string;
  recordedOn: Date;
  episodes?: Episode[];
  movieYear?: number;
}

export interface Episode {
  id?: number;
  recordedOn?: Date;
  firstAiredOn?: Date;
  season?: number;
  episode?: number;
  episodeTitle?: string;
}


export function ShowListing() {
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(() => {
    fetch("input/shows.json")
      .then((response) => response.json())
      .then((data) => {
        const parsedShows = data.map((show: Show) => ({
          ...show,
          recordedOn: new Date(show.recordedOn),
        }));
        setShows(parsedShows);
      })
      .catch((error) => console.error(error));
  }, []);

  return (
    <div>
      <label style={{ fontSize: "2rem", color: "lightblue", paddingLeft: "1rem" }}>My Shows</label>
      <table className="showListingTable" style={{ marginTop: "30px" }}>
        <tbody>
          {shows.map((show: Show, index: number) => {
            const episodeCount = show.episodes?.length || 0
            const episodeCountLabel = episodeCount > 1 ? `[${episodeCount}]` : ""
            const icon = show.movieYear ? "./images/movie.png" :
              episodeCount === 1 ? "./images/television.png" : "./images/folder.png"
            const dayOfWeek = show.recordedOn.toLocaleDateString('en-US', { weekday: 'short' });
            const monthDay = show.recordedOn.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit' });

            return (
              <tr key={index}>
                <td>
                  <img
                    src={icon}
                    style={{ width: "3rem", height: "3rem" }}
                    alt=""
                  />
                </td>
                <td>{show.title} {episodeCountLabel}</td>
                <td>{dayOfWeek}</td>
                <td>{monthDay}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
