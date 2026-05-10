import { Show, ShowKind, Series, Episode, Movie } from "../services/generated/graphql-types"

const toDate = (value: any): Date | undefined => {
  if (!value) return undefined;
  if (value instanceof Date) return value;
  if (typeof value === 'string') return new Date(value);
  return undefined;
};

export const recordedOn = (show: Show): Date | undefined => {
  const movieDate = toDate((show as Movie).recordedOn);
  if (movieDate) return movieDate;
  
  const episodeDate = toDate((show as Episode).recordedOn);
  if (episodeDate) return episodeDate;
  
  const series = show as Series;
  if (series.episodes && series.episodes.length > 0) {
    return series.episodes.reduce(
      (latest, episode) => {
        const episodeDate = toDate(episode.recordedOn);
        return episodeDate && episodeDate > latest ? episodeDate : latest;
      },
      new Date(0)
    );
  }
  
  return undefined;
};

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

export const mergeEpisodes = (shows: Show[]): Show[] => {
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
