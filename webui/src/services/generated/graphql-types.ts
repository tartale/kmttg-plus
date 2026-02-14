export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Any: { input: any; output: any; }
  Date: { input: string; output: Date; }
  Time: { input: string; output: Date; }
};

export type Episode = Show & {
  __typename?: 'Episode';
  description: Scalars['String']['output'];
  episodeDescription: Scalars['String']['output'];
  episodeNumber: Scalars['Int']['output'];
  episodeTitle: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  kind: ShowKind;
  originalAirDate: Scalars['Date']['output'];
  recordedOn: Scalars['Time']['output'];
  seasonNumber: Scalars['Int']['output'];
  seriesId: Scalars['ID']['output'];
  title: Scalars['String']['output'];
};

export type EpisodeFilter = {
  description?: InputMaybe<FilterOperator>;
  episodeDescription?: InputMaybe<FilterOperator>;
  episodeNumber?: InputMaybe<FilterOperator>;
  episodeTitle?: InputMaybe<FilterOperator>;
  kind?: InputMaybe<FilterOperator>;
  originalAirDate?: InputMaybe<FilterOperator>;
  recordedOn?: InputMaybe<FilterOperator>;
  seasonNumber?: InputMaybe<FilterOperator>;
  title?: InputMaybe<FilterOperator>;
};

export type FilterOperator = {
  eq?: InputMaybe<Scalars['Any']['input']>;
  gt?: InputMaybe<Scalars['Any']['input']>;
  gte?: InputMaybe<Scalars['Any']['input']>;
  lt?: InputMaybe<Scalars['Any']['input']>;
  lte?: InputMaybe<Scalars['Any']['input']>;
  matches?: InputMaybe<Scalars['Any']['input']>;
  ne?: InputMaybe<Scalars['Any']['input']>;
};

export type Movie = Show & {
  __typename?: 'Movie';
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  imageURL: Scalars['String']['output'];
  kind: ShowKind;
  movieYear: Scalars['Int']['output'];
  recordedOn: Scalars['Time']['output'];
  title: Scalars['String']['output'];
};


export type MovieImageUrlArgs = {
  height?: InputMaybe<Scalars['Int']['input']>;
  width?: InputMaybe<Scalars['Int']['input']>;
};

export type MovieFilter = {
  description?: InputMaybe<FilterOperator>;
  kind?: InputMaybe<FilterOperator>;
  movieYear?: InputMaybe<FilterOperator>;
  recordedOn?: InputMaybe<FilterOperator>;
  title?: InputMaybe<FilterOperator>;
};

export type Query = {
  __typename?: 'Query';
  tivos: Array<Tivo>;
};


export type QueryTivosArgs = {
  filters?: InputMaybe<Array<InputMaybe<TivoFilter>>>;
};

export type Series = Show & {
  __typename?: 'Series';
  description: Scalars['String']['output'];
  episodes: Array<Episode>;
  id: Scalars['ID']['output'];
  imageURL: Scalars['String']['output'];
  kind: ShowKind;
  recordedOn: Scalars['Time']['output'];
  title: Scalars['String']['output'];
};


export type SeriesEpisodesArgs = {
  filter?: InputMaybe<Array<InputMaybe<EpisodeFilter>>>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type SeriesImageUrlArgs = {
  height?: InputMaybe<Scalars['Int']['input']>;
  width?: InputMaybe<Scalars['Int']['input']>;
};

export type SeriesFilter = {
  description?: InputMaybe<FilterOperator>;
  kind?: InputMaybe<FilterOperator>;
  recordedOn?: InputMaybe<FilterOperator>;
  title?: InputMaybe<FilterOperator>;
};

export type Show = {
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  kind: ShowKind;
  recordedOn: Scalars['Time']['output'];
  title: Scalars['String']['output'];
};

export type ShowFilter = {
  and?: InputMaybe<Array<InputMaybe<ShowFilter>>>;
  description?: InputMaybe<FilterOperator>;
  kind?: InputMaybe<FilterOperator>;
  or?: InputMaybe<Array<InputMaybe<ShowFilter>>>;
  recordedOn?: InputMaybe<FilterOperator>;
  title?: InputMaybe<FilterOperator>;
};

export enum ShowKind {
  Episode = 'EPISODE',
  Movie = 'MOVIE',
  Series = 'SERIES'
}

export type SortBy = {
  direction: SortDirection;
  field: Scalars['Any']['input'];
};

export enum SortDirection {
  Asc = 'ASC',
  Desc = 'DESC'
}

export type Sorter = {
  fields: Array<SortBy>;
};

export type Tivo = {
  __typename?: 'Tivo';
  address: Scalars['String']['output'];
  name: Scalars['String']['output'];
  shows?: Maybe<Array<Maybe<Show>>>;
  tsn: Scalars['String']['output'];
};


export type TivoShowsArgs = {
  filters?: InputMaybe<Array<InputMaybe<ShowFilter>>>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};

export type TivoFilter = {
  name?: InputMaybe<FilterOperator>;
};
