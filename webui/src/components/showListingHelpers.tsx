import { Show, ShowKind, Series, Episode, Movie } from "../services/generated/graphql-types"

type ShowSetter = React.Dispatch<React.SetStateAction<Show[]>>

export const getShows = (setShows: ShowSetter) => () => {
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
      const shows = mergeEpisodes(jsonArray.map(parseShow))
      setShows(shows);
    })
    .catch((error) => console.error(error));

}

export const recordedOn = (show: Show): Date | undefined =>
  (show as Movie).recordedOn ||
  (show as Episode).recordedOn ||
  (show as Series).episodes.reduce(
    (latest, episode) =>
      episode.recordedOn > latest ? episode.recordedOn : latest,
    new Date(0)
  ) ||
  undefined;

export const getImageFileForShow = (show: Show, open: boolean): string => {
  switch (show.kind) {
    case ShowKind.Series: {
      if (open) {
        return "./images/folder-open.png";
      } else {
        return "./images/folder-closed.png";
      }
    }
    case ShowKind.Episode: {
      return "./images/television.png";
    }
    case ShowKind.Movie: {
      return "./images/movie.png";
    }
    default: {
      return "./images/television-unknown.png";
    }
  }
};

export const getTitleExtension = (show: Show): string => {
  var titleExtension = ""
  switch (show.kind) {
    case ShowKind.Movie:
      const movie = (show as Movie);
      titleExtension = `(${movie.movieYear})`
      break
    case ShowKind.Series:
      const series = (show as Series);
      const episodeCount = series.episodes.length
      titleExtension = `[${episodeCount}]`
      break
    case ShowKind.Episode:
      const episode = (show as Episode);
      const seasonLabel = episode.seasonNumber ? `S${episode.seasonNumber.toString().padStart(2, '0')}` : ""
      const episodeLabel = episode.episodeNumber ? `E${episode.episodeNumber.toString().padStart(2, '0')}` : ""
      titleExtension = `[${seasonLabel}${episodeLabel}]`
      break
  }

  return titleExtension
}

const parseShow = (obj: any): Show => {
  const recording = obj.recording[0]
  const show: Series | Movie = {
    recordingID: recording.recordingID,
    kind: recording.episodic ? ShowKind.Series : ShowKind.Movie,
    title: recording.title,
    recordedOn: new Date(recording.startTime),
    description: recording.description,
    movieYear: recording.movieYear,
    episodes: recording.episodic ? [
      {
        recordingID: recording.recordingID,
        kind: ShowKind.Episode,
        title: recording.title,
        recordedOn: new Date(recording.startTime),
        description: recording.description,
        originalAirDate: new Date(recording.originalAirDate),
        seasonNumber: recording.seasonNumber ? recording.seasonNumber : 0,
        episodeNumber: recording.episodeNum ? recording.episodeNum[0] : 0,
        episodeTitle: recording.subtitle,
        episodeDescription: recording.description,
      }
    ] : []
  }

  return show
}

const mergeEpisodes = (shows: Show[]): Show[] => {
  const combinedShows: Show[] = Object.values(shows.reduce((acc: any, show) => {
    if (show.kind === ShowKind.Series) {
      const series = (show as Series)
      if (acc[series.title]) {
        acc[series.title].episodes.push(...series.episodes);
      } else {
        acc[series.title] = { ...series };
      }
    } else {
      acc[show.title] = { ...show }
    }
    return acc;
  }, {}));

  return combinedShows
}
