import json
from wsgiref.simple_server import make_server
from pyramid.config import Configurator
from pyramid.view import view_config
from pyramid.response import Response
from pyramid.renderers import render_to_response


import read


@view_config(route_name='upload', request_method='POST')
def upload(request):
    if 'image' not in request.POST:
        return Response('Arquivo de imagem não encontrado', status_code=400)
    image_file = request.POST['image'].file
    read.reader(image_file)
    return Response('Imagem recebida com sucesso', status_code=200)


@view_config(route_name='upload_images', request_method='POST')
def upload_images(request):
    images = {}
    files = request.POST.getall('file') # obter a lista de arquivos enviados
    num_expected_images = len(files)  # definir o número de imagens esperadas
    for file in files:
        # processar o arquivo como antes
        images[file.name] = file.file.read()  # adiciona a imagem ao dicionário de imagens
    print(len(images))
    if len(images) == num_expected_images:
        result = {'result': 'Imagens Recebidas', 'imagens': images}
        print(result)
        return Response(json.dumps(result), content_type='application/json', status_code=200, charset='utf-8')
    else:
        result = {'error': 'Número de imagens recebidas não é igual ao número esperado.'}
        return Response(json.dumps(result), content_type='application/json', status_code=400, charset='utf-8')




if __name__ == '__main__':
    config = Configurator()
    config.add_route('upload', '/upload_image')
    config.add_route('upload_images', '/upload_images')
    config.scan()
    app = config.make_wsgi_app()
    server = make_server('0.0.0.0', 8777, app)
    server.serve_forever()
