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
  Date: { input: string; output: Date; }
  Time: { input: string; output: Date; }
};

export type Episode = Show & {
  __typename?: 'Episode';
  description: Scalars['String']['output'];
  episodeDescription: Scalars['String']['output'];
  episodeNumber: Scalars['Int']['output'];
  episodeTitle: Scalars['String']['output'];
  kind: ShowKind;
  originalAirDate: Scalars['Date']['output'];
  recordedOn: Scalars['Time']['output'];
  recordingID: Scalars['String']['output'];
  seasonNumber: Scalars['Int']['output'];
  title: Scalars['String']['output'];
};

export type Movie = Show & {
  __typename?: 'Movie';
  description: Scalars['String']['output'];
  kind: ShowKind;
  movieYear?: Maybe<Scalars['String']['output']>;
  recordedOn: Scalars['Time']['output'];
  recordingID: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  tivos: Array<Tivo>;
};

export type Series = Show & {
  __typename?: 'Series';
  description: Scalars['String']['output'];
  episodes: Array<Episode>;
  kind: ShowKind;
  recordedOn: Scalars['Time']['output'];
  recordingID: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type Show = {
  description: Scalars['String']['output'];
  kind: ShowKind;
  recordedOn: Scalars['Time']['output'];
  recordingID: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export enum ShowKind {
  Episode = 'EPISODE',
  Movie = 'MOVIE',
  Series = 'SERIES'
}

export type Tivo = {
  __typename?: 'Tivo';
  address: Scalars['String']['output'];
  name: Scalars['String']['output'];
  recordings?: Maybe<Array<Maybe<Show>>>;
  tsn: Scalars['String']['output'];
};


export type TivoRecordingsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};
