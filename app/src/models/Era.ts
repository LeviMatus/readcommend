export type Era = {
  id: number;
  title: string;
  minYear?: number;
  maxYear?: number;
};

export const DefaultEra: Era = {
  id: 0,
  title: "",
};
