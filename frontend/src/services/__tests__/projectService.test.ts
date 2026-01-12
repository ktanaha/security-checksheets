import axios from 'axios';
import {
  getProjects,
  getProject,
  createProject,
  updateProject,
  deleteProject,
  getProjectFiles,
  uploadFile,
  deleteFile,
} from '../projectService';
import type { Project, CreateProjectRequest, UpdateProjectRequest, UploadedFile } from '../../types/project';

// axiosのモック
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('projectService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getProjects', () => {
    it('案件一覧を取得できる', async () => {
      // Arrange
      const mockProjects: Project[] = [
        {
          id: 1,
          customer_name: '株式会社テスト',
          description: 'テスト案件',
          owner: 'テスト太郎',
          status: 'active',
          created_at: '2026-01-12T00:00:00Z',
          updated_at: '2026-01-12T00:00:00Z',
        },
      ];

      mockedAxios.get.mockResolvedValueOnce({ data: mockProjects });

      // Act
      const result = await getProjects();

      // Assert
      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:8080/api/projects');
      expect(result).toEqual(mockProjects);
    });

    it('エラー時に例外をスローする', async () => {
      // Arrange
      mockedAxios.get.mockRejectedValueOnce(new Error('Network error'));

      // Act & Assert
      await expect(getProjects()).rejects.toThrow('Network error');
    });
  });

  describe('getProject', () => {
    it('案件詳細を取得できる', async () => {
      // Arrange
      const mockProject: Project = {
        id: 1,
        customer_name: '株式会社テスト',
        description: 'テスト案件',
        owner: 'テスト太郎',
        status: 'active',
        created_at: '2026-01-12T00:00:00Z',
        updated_at: '2026-01-12T00:00:00Z',
      };

      mockedAxios.get.mockResolvedValueOnce({ data: mockProject });

      // Act
      const result = await getProject(1);

      // Assert
      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:8080/api/projects/1');
      expect(result).toEqual(mockProject);
    });
  });

  describe('createProject', () => {
    it('案件を作成できる', async () => {
      // Arrange
      const request: CreateProjectRequest = {
        customer_name: '株式会社テスト',
        description: 'テスト案件',
        owner: 'テスト太郎',
      };

      const mockResponse: Project = {
        id: 1,
        ...request,
        status: 'active',
        created_at: '2026-01-12T00:00:00Z',
        updated_at: '2026-01-12T00:00:00Z',
      };

      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await createProject(request);

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8080/api/projects',
        request
      );
      expect(result).toEqual(mockResponse);
    });
  });

  describe('updateProject', () => {
    it('案件を更新できる', async () => {
      // Arrange
      const request: UpdateProjectRequest = {
        customer_name: '株式会社テスト（更新）',
        description: 'テスト案件（更新）',
        owner: 'テスト花子',
        status: 'in_progress',
      };

      const mockResponse: Project = {
        id: 1,
        ...request,
        created_at: '2026-01-12T00:00:00Z',
        updated_at: '2026-01-12T01:00:00Z',
      };

      mockedAxios.put.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await updateProject(1, request);

      // Assert
      expect(mockedAxios.put).toHaveBeenCalledWith(
        'http://localhost:8080/api/projects/1',
        request
      );
      expect(result).toEqual(mockResponse);
    });
  });

  describe('deleteProject', () => {
    it('案件を削除できる', async () => {
      // Arrange
      mockedAxios.delete.mockResolvedValueOnce({ data: null });

      // Act
      await deleteProject(1);

      // Assert
      expect(mockedAxios.delete).toHaveBeenCalledWith('http://localhost:8080/api/projects/1');
    });
  });

  describe('getProjectFiles', () => {
    it('案件のファイル一覧を取得できる', async () => {
      // Arrange
      const mockFiles: UploadedFile[] = [
        {
          id: 1,
          project_id: 1,
          file_name: 'test.xlsx',
          file_path: '/app/uploads/project_1/test.xlsx',
          file_size: 12345,
          uploaded_by: 'テスト太郎',
          uploaded_at: '2026-01-12T00:00:00Z',
        },
      ];

      mockedAxios.get.mockResolvedValueOnce({ data: mockFiles });

      // Act
      const result = await getProjectFiles(1);

      // Assert
      expect(mockedAxios.get).toHaveBeenCalledWith(
        'http://localhost:8080/api/projects/1/files'
      );
      expect(result).toEqual(mockFiles);
    });
  });

  describe('uploadFile', () => {
    it('ファイルをアップロードできる', async () => {
      // Arrange
      const file = new File(['test'], 'test.xlsx', {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      });
      const uploadedBy = 'テスト太郎';

      const mockResponse: UploadedFile = {
        id: 1,
        project_id: 1,
        file_name: 'test.xlsx',
        file_path: '/app/uploads/project_1/test.xlsx',
        file_size: 4,
        uploaded_by: uploadedBy,
        uploaded_at: '2026-01-12T00:00:00Z',
      };

      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      // Act
      const result = await uploadFile(1, { file, uploaded_by: uploadedBy });

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8080/api/projects/1/files',
        expect.any(FormData),
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      );
      expect(result).toEqual(mockResponse);
    });
  });

  describe('deleteFile', () => {
    it('ファイルを削除できる', async () => {
      // Arrange
      mockedAxios.delete.mockResolvedValueOnce({ data: null });

      // Act
      await deleteFile(1);

      // Assert
      expect(mockedAxios.delete).toHaveBeenCalledWith('http://localhost:8080/api/files/1');
    });
  });
});
