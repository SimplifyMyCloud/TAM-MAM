// src/App.js
import React, { useState, useEffect } from 'react';
import { Search, Video, FileAudio, Grid, List, Upload } from 'lucide-react';
import UploadDialog from './components/UploadDialog';

function App() {
  const [viewMode, setViewMode] = useState('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [assets, setAssets] = useState([]);
  const [isUploadOpen, setIsUploadOpen] = useState(false);

  // Mock data
  const mockAssets = [
    {
      id: "1",
      title: "Brand Video 2024",
      type: "video",
      thumbnail: "/api/placeholder/320/180",
      metadata: {
        duration: "2:30",
        format: "MP4",
        resolution: "1920x1080"
      }
    },
    {
      id: "2",
      title: "Company Podcast",
      type: "audio",
      thumbnail: "/api/placeholder/320/180",
      metadata: {
        duration: "25:00",
        format: "MP3",
        bitrate: "320kbps"
      }
    }
  ];

  useEffect(() => {
    setAssets(mockAssets);
  }, []);

  return (
    <div className="max-w-7xl mx-auto p-4 space-y-4">
      {/* Header */}
      <div className="bg-white rounded-lg shadow p-6">
        <h1 className="text-2xl font-bold">Media Asset Manager</h1>
        <p className="text-gray-500">Manage and organize your media assets</p>
      </div>

      {/* Search and Controls */}
      <div className="flex items-center justify-between gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-3 h-4 w-4 text-gray-500" />
          <input
            type="text"
            placeholder="Search assets..."
            className="w-full p-2 pl-10 border rounded-lg"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <div className="flex items-center gap-2">
          <button 
            className={`p-2 rounded-lg ${viewMode === 'grid' ? 'bg-blue-500 text-white' : 'bg-gray-100'}`}
            onClick={() => setViewMode('grid')}
          >
            <Grid className="h-4 w-4" />
          </button>
          <button
            className={`p-2 rounded-lg ${viewMode === 'list' ? 'bg-blue-500 text-white' : 'bg-gray-100'}`}
            onClick={() => setViewMode('list')}
          >
            <List className="h-4 w-4" />
          </button>
          <button 
            className="flex items-center gap-2 bg-blue-500 text-white px-4 py-2 rounded-lg"
            onClick={() => setIsUploadOpen(true)}
          >
            <Upload className="h-4 w-4" />
            Upload
          </button>
        </div>
      </div>

      {/* Grid View */}
      {viewMode === 'grid' ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {assets.map(asset => (
            <div key={asset.id} className="bg-white rounded-lg shadow overflow-hidden">
              <div className="relative aspect-video">
                <img 
                  src={asset.thumbnail} 
                  alt={asset.title}
                  className="object-cover w-full h-full"
                />
                {asset.type === 'video' && (
                  <div className="absolute bottom-2 right-2 bg-black/50 text-white px-2 py-1 rounded text-sm">
                    {asset.metadata.duration}
                  </div>
                )}
              </div>
              <div className="p-4">
                <h3 className="font-medium truncate">{asset.title}</h3>
                <div className="text-sm text-gray-500 space-y-1">
                  <div className="flex items-center gap-2">
                    {asset.type === 'video' ? (
                      <Video className="h-4 w-4" />
                    ) : (
                      <FileAudio className="h-4 w-4" />
                    )}
                    <span>{asset.type}</span>
                  </div>
                  <div className="text-xs space-x-2">
                    <span>{asset.metadata.format}</span>
                    {asset.metadata.resolution && (
                      <span>• {asset.metadata.resolution}</span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        // List View
        <div className="space-y-2">
          {assets.map(asset => (
            <div key={asset.id} className="bg-white rounded-lg shadow">
              <div className="p-4 flex items-center gap-4">
                <div className="w-40 aspect-video relative flex-shrink-0">
                  <img 
                    src={asset.thumbnail} 
                    alt={asset.title}
                    className="object-cover w-full h-full rounded"
                  />
                </div>
                <div className="flex-grow">
                  <h3 className="font-medium">{asset.title}</h3>
                  <div className="text-sm text-gray-500 space-y-1">
                    <div className="flex items-center gap-2">
                      {asset.type === 'video' ? (
                        <Video className="h-4 w-4" />
                      ) : (
                        <FileAudio className="h-4 w-4" />
                      )}
                      <span>{asset.type}</span>
                    </div>
                    <div className="space-x-2">
                      <span>{asset.metadata.format}</span>
                      {asset.metadata.resolution && (
                        <span>• {asset.metadata.resolution}</span>
                      )}
                      <span>• {asset.metadata.duration}</span>
                    </div>
                  </div>
                </div>
                <div className="flex-shrink-0">
                  <button className="px-4 py-2 border rounded-lg hover:bg-gray-50">
                    View Details
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Upload Dialog */}
      <UploadDialog 
        isOpen={isUploadOpen} 
        onClose={() => setIsUploadOpen(false)} 
      />
    </div>
  );
}

export default App;