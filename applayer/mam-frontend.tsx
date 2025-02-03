import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, Video, FileAudio, Database, Layers, Grid, List, Upload } from 'lucide-react';

// Main App Component
const MAMApp = () => {
  const [view, setView] = useState('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [assets, setAssets] = useState([]);
  const [loading, setLoading] = useState(false);

  // Mock data - replace with API calls
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
    <div className="w-full max-w-7xl mx-auto p-4 space-y-4">
      {/* Header */}
      <Card>
        <CardHeader>
          <CardTitle>Media Asset Manager</CardTitle>
          <CardDescription>Manage and organize your media assets</CardDescription>
        </CardHeader>
      </Card>

      {/* Search and Controls */}
      <div className="flex items-center justify-between gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            placeholder="Search assets..."
            className="pl-8"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <div className="flex items-center gap-2">
          <Button 
            variant={view === 'grid' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setView('grid')}
          >
            <Grid className="h-4 w-4" />
          </Button>
          <Button
            variant={view === 'list' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setView('list')}
          >
            <List className="h-4 w-4" />
          </Button>
          <Button>
            <Upload className="h-4 w-4 mr-2" />
            Upload
          </Button>
        </div>
      </div>

      {/* Main Content */}
      {view === 'grid' ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {assets.map(asset => (
            <AssetCard key={asset.id} asset={asset} />
          ))}
        </div>
      ) : (
        <div className="space-y-2">
          {assets.map(asset => (
            <AssetListItem key={asset.id} asset={asset} />
          ))}
        </div>
      )}

      {/* Upload Dialog */}
      <UploadDialog />
    </div>
  );
};

// Asset Card Component
const AssetCard = ({ asset }) => {
  return (
    <Card className="overflow-hidden">
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
      <CardContent className="p-4">
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
      </CardContent>
    </Card>
  );
};

// Asset List Item Component
const AssetListItem = ({ asset }) => {
  return (
    <Card>
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
          <Button variant="outline" size="sm">View Details</Button>
        </div>
      </div>
    </Card>
  );
};

// Upload Dialog Component
const UploadDialog = () => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleUpload = () => {
    setUploading(true);
    // Simulate upload progress
    const interval = setInterval(() => {
      setProgress(prev => {
        if (prev >= 100) {
          clearInterval(interval);
          setUploading(false);
          return 0;
        }
        return prev + 10;
      });
    }, 500);
  };

  return (
    <Card className="p-4">
      <div className="space-y-4">
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
          <Upload className="h-8 w-8 mx-auto text-gray-400 mb-4" />
          <p className="text-sm text-gray-500">
            Drag and drop files here, or click to select files
          </p>
        </div>
        {uploading && (
          <div className="w-full bg-gray-200 rounded-full h-2.5">
            <div 
              className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
        )}
        <Button onClick={handleUpload} disabled={uploading}>
          {uploading ? 'Uploading...' : 'Upload Files'}
        </Button>
      </div>
    </Card>
  );
};

export default MAMApp;
