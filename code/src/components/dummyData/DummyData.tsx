type ContentItem = {
  id : string;
  type: string;
  url: string;
  title: string;
  size: string;
};

interface PostData {
  id: string,
  title: string,
  content: ContentItem[],
  description: string,
  membership: string[],
  createdAt: string,
}

const Posts: PostData[] = [
  { 
    id: '1', 
    title: 'First Post', 
    content: [
    {
      id: '3',
      type: 'image',
      url: 'https://utfs.io/f/...',
      title:'image.png',
      size: '1.2 MB'
    },
    {
      id: '4',
      type: 'video',
      url: 'https://utfs.io/f/...',
      title: 'video.mp4',
      size: '5.4 MB',
    },
    {
      id: '5',
      type: 'document',
      url: 'https://utfs.io/f/...',
      title: 'Project Notes.pdf',
      size: '2.4 MB'
    }
  ],
    description: "This is your first post. You can edit it or create new posts to share with your subscribers. The description provides a preview of your post content.", 
    membership: ['Gold Tier'],
    createdAt: 'June 15, 2025'
  },
];

export default Posts