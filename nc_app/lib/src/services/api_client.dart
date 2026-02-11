import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiClient {
  ApiClient(this.baseUrl, {this.token});

  final String baseUrl;
  final String? token;

  /// Для Android-эмулятора: http://10.0.2.2:8080
  /// Для реального устройства: IP компьютера, например http://192.168.1.100:8080
  static const defaultBaseUrl = 'http://localhost:8080';

  Map<String, String> _headers({bool jsonBody = true}) {
    final headers = <String, String>{};
    if (jsonBody) {
      headers['Content-Type'] = 'application/json; charset=utf-8';
    }
    if (token != null && token!.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }
    return headers;
  }

  Future<Map<String, dynamic>> postJson(String path, Map<String, dynamic> body) async {
    final uri = Uri.parse('$baseUrl$path');
    final resp = await http.post(uri, headers: _headers(), body: jsonEncode(body));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      return jsonDecode(utf8.decode(resp.bodyBytes)) as Map<String, dynamic>;
    }
    throw ApiException(resp.statusCode, resp.body);
  }

  Future<List<dynamic>> getList(String path) async {
    final uri = Uri.parse('$baseUrl$path');
    final resp = await http.get(uri, headers: _headers(jsonBody: false));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      final decoded = jsonDecode(utf8.decode(resp.bodyBytes));
      return decoded is List ? decoded : <dynamic>[];
    }
    throw ApiException(resp.statusCode, resp.body);
  }

  Future<Map<String, dynamic>> getJson(String path) async {
    final uri = Uri.parse('$baseUrl$path');
    final resp = await http.get(uri, headers: _headers(jsonBody: false));
    if (resp.statusCode >= 200 && resp.statusCode < 300) {
      return jsonDecode(utf8.decode(resp.bodyBytes)) as Map<String, dynamic>;
    }
    throw ApiException(resp.statusCode, resp.body);
  }
}

class ApiException implements Exception {
  ApiException(this.statusCode, this.body);
  final int statusCode;
  final String body;

  @override
  String toString() => 'ApiException($statusCode, $body)';
}
