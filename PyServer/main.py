from wsgiref.simple_server import make_server
from pyramid.config import Configurator
from pyramid.view import view_config
from pyramid.response import Response

import read


@view_config(route_name='upload', request_method='POST')
def upload(request):
    if 'image' not in request.POST:
        return Response('Arquivo de imagem não encontrado', status_code=400)
    image_file = request.POST['image'].file
    read.reader(image_file)
    return Response('Imagem recebida com sucesso', status_code=200)



if __name__ == '__main__':
    config = Configurator()
    config.add_route('upload', '/upload_image')
    config.add_route('uploads', '/upload_images')
    config.scan()
    app = config.make_wsgi_app()
    server = make_server('0.0.0.0', 8777, app)
    server.serve_forever()
