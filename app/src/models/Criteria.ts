export type Criteria = {
  authors: number[];
  genres: number[];
  minYear?: number;
  maxYear?: number;
  minPages?: number;
  maxPages?: number;
  limit?: number;
};

export const DefaultCriteria: Criteria = {
  authors: [],
  genres: [],
  limit: 10,
};
