import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, Video, FileAudio, Database, Layers } from 'lucide-react';

const TAMSDashboard = () => {
  const [activeTab, setActiveTab] = useState('flows');
  const [flows, setFlows] = useState([]);
  const [sources, setSources] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(false);

  // Mock data for demonstration
  const mockFlows = [
    {
      id: "550e8400-e29b-41d4-a716-446655440000",
      label: "Main Camera Feed",
      format: "urn:x-nmos:format:video",
      codec: "video/h264",
      tags: { location: "studio-1" }
    },
    {
      id: "550e8400-e29b-41d4-a716-446655440001",
      label: "Audio Feed 1",
      format: "urn:x-nmos:format:audio",
      codec: "audio/aac",
      tags: { channel: "stereo" }
    }
  ];

  const mockSources = [
    {
      id: "550e8400-e29b-41d4-a716-446655440002",
      label: "Camera 1",
      format: "urn:x-nmos:format:video",
      description: "Main studio camera"
    },
    {
      id: "550e8400-e29b-41d4-a716-446655440003",
      label: "Microphone Set 1",
      format: "urn:x-nmos:format:audio",
      description: "Studio microphone array"
    }
  ];

  useEffect(() => {
    // Simulating API fetch
    setFlows(mockFlows);
    setSources(mockSources);
  }, []);

  const getFormatIcon = (format) => {
    switch (format) {
      case 'urn:x-nmos:format:video':
        return <Video className="w-4 h-4" />;
      case 'urn:x-nmos:format:audio':
        return <FileAudio className="w-4 h-4" />;
      case 'urn:x-nmos:format:data':
        return <Database className="w-4 h-4" />;
      case 'urn:x-nmos:format:multi':
        return <Layers className="w-4 h-4" />;
      default:
        return null;
    }
  };

  return (
    <div className="w-full max-w-6xl mx-auto p-4 space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Time-addressable Media Store</CardTitle>
          <CardDescription>
            Manage and monitor segmented media flows
          </CardDescription>
        </CardHeader>
      </Card>

      <div className="flex items-center space-x-2">
        <div className="relative flex-1">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            placeholder="Search..."
            className="pl-8"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <Button>Add New</Button>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="flows">Flows</TabsTrigger>
          <TabsTrigger value="sources">Sources</TabsTrigger>
        </TabsList>

        <TabsContent value="flows">
          <Card>
            <CardHeader>
              <CardTitle>Media Flows</CardTitle>
              <CardDescription>List of all registered media flows</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative overflow-x-auto">
                <table className="w-full text-sm text-left">
                  <thead className="text-xs uppercase bg-gray-50">
                    <tr>
                      <th className="px-6 py-3">Format</th>
                      <th className="px-6 py-3">Label</th>
                      <th className="px-6 py-3">Codec</th>
                      <th className="px-6 py-3">ID</th>
                      <th className="px-6 py-3">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {flows.map((flow) => (
                      <tr key={flow.id} className="bg-white border-b hover:bg-gray-50">
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            {getFormatIcon(flow.format)}
                            {flow.format.split(':').pop()}
                          </div>
                        </td>
                        <td className="px-6 py-4">{flow.label}</td>
                        <td className="px-6 py-4">{flow.codec}</td>
                        <td className="px-6 py-4 font-mono text-sm">{flow.id}</td>
                        <td className="px-6 py-4">
                          <Button variant="outline" size="sm">View Details</Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="sources">
          <Card>
            <CardHeader>
              <CardTitle>Media Sources</CardTitle>
              <CardDescription>List of all registered media sources</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative overflow-x-auto">
                <table className="w-full text-sm text-left">
                  <thead className="text-xs uppercase bg-gray-50">
                    <tr>
                      <th className="px-6 py-3">Format</th>
                      <th className="px-6 py-3">Label</th>
                      <th className="px-6 py-3">Description</th>
                      <th className="px-6 py-3">ID</th>
                      <th className="px-6 py-3">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {sources.map((source) => (
                      <tr key={source.id} className="bg-white border-b hover:bg-gray-50">
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            {getFormatIcon(source.format)}
                            {source.format.split(':').pop()}
                          </div>
                        </td>
                        <td className="px-6 py-4">{source.label}</td>
                        <td className="px-6 py-4">{source.description}</td>
                        <td className="px-6 py-4 font-mono text-sm">{source.id}</td>
                        <td className="px-6 py-4">
                          <Button variant="outline" size="sm">View Details</Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default TAMSDashboard;
