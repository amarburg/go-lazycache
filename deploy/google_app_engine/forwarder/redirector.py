import webapp2

def get_redirect_uri(handler, *args, **kwargs):
    return 'https://camhd-cache.appspot.com/' + kwargs.get('path')

app = webapp2.WSGIApplication([
    webapp2.Route('/<path:.*>', webapp2.RedirectHandler, defaults={'_uri': get_redirect_uri}),
], debug=False)
