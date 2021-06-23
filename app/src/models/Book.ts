export type Book = {
  id: number;
  title: string;
  yearPublished: number;
  rating: number;
  pages: number;
  genre: {
    id: number;
    title: string;
  };
  author: {
    id: number;
    firstName: string;
    lastName: string;
  };
};
