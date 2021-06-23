export type Size = {
  id: number;
  title: string;
  minPages?: number;
  maxPages?: number;
};

export const DefaultSize = {
  id: 0,
  title: "",
};
