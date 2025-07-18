import ContentRenderer from "./ContentRenderer";

export default function CreatorHome() {
  const posts = [
    { 
      id: 1, 
      title: 'First Post', 
      content: [
      {
        type: 'image',
        url: 'https://utfs.io/f/...',
        title:'image.png',
        size: '1.2 MB'
      },
      {
        type: 'video',
        url: 'https://utfs.io/f/...',
        title: 'video.mp4',
        size: '5.4 MB',
      },
      {
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

  return (
    <div className="w-full p-4">
      {/* Header with title and add button */}
      <div className="mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h2 className="text-2xl font-bold text-white">Your Posts</h2>
        
        <button className="px-4 py-2 text-sm bg-white text-black rounded-lg hover:bg-gray-200 flex items-center justify-center transition-colors cursor-pointer">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
          </svg>
          Add New Post
        </button>
      </div>

      {/* Posts list */}
      <div className="space-y-4">
        {posts.map((post) => (
          <div 
            key={post.id}
            className="bg-neutral-800 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 border border-neutral-700"
          >
            <div className="p-4">
              <div className="flex justify-between items-start mb-2">
                <h3 className="text-lg font-semibold text-white">{post.title}</h3>
                <span className="text-neutral-400 text-sm">{post.createdAt}</span>
              </div>

              {/* Post content  */}
              {post.content && post.content.length > 0 && (
                <ContentRenderer contentItem={post.content} />
              )}
              
              {post.description && (
                <p className="mt-2 text-neutral-300 text-sm">{post.description}</p>
              )}
              
              <div className="mt-4 flex justify-end">
                <button className="px-3 py-1 text-sm bg-white text-black rounded-md hover:bg-gray-200 transition-colors cursor-pointer">
                  Edit Post
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Empty state (when no posts) */}
      {posts.length === 0 && (
        <div className="bg-neutral-800 rounded-lg p-8 text-center border border-neutral-700">
          <p className="text-neutral-400 mb-4">You haven't created any posts yet</p>
          <button className="px-4 py-2 bg-white text-black rounded-lg hover:bg-gray-200 inline-flex items-center transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
            </svg>
            Create Your First Post
          </button>
        </div>
      )}
    </div>
  );
}