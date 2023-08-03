import { Show, ShowKind, Series, Episode, Movie } from "../services/generated/graphql-types"

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
